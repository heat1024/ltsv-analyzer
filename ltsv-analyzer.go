package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	bufferSize = 4096
)

type result struct {
	sum   uint64
	avg   uint64
	count uint64
}

type sortedKey struct {
	Key   string
	Value uint64
}

func ltsvAnalyze(tFile []string) (map[string]*result, []sortedKey, error) {
	results := make(map[string]*result)
	var resultKeys []sortedKey
	var isCompressed bool

	if len(tFiles) == 0 {
		tFiles = append(tFiles, "STDIN")
	}

	for _, file := range tFiles {
		var logFile *os.File
		var r *bufio.Reader
		var zr *gzip.Reader
		var err error

		logFile = nil

		if file == "STDIN" {
			isCompressed = false
			r = bufio.NewReader(os.Stdin)
		} else {
			logFile, err = os.OpenFile(file, os.O_RDONLY, 0400)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot open log file %s: %s", file, err.Error())
			}
			isCompressed = false

			// read first 512 bytes for detect contents type.
			// why 512 bytes? ref : https://golang.org/pkg/net/http/#DetectContentType
			buf := make([]byte, 512)
			_, err = logFile.Read(buf)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot read file header: %s", err.Error())
			}
			logFile.Seek(0, 0)

			filetype := http.DetectContentType(buf)

			switch filetype {
			case "application/x-gzip":
				zr, err = gzip.NewReader(logFile)
				if err != nil {
					return nil, nil, fmt.Errorf("cannot create gzip reader for %s: %s", file, err.Error())
				}
				isCompressed = true
			default:
				isCompressed = false
				r = bufio.NewReader(logFile)
			}
		}

		var leftLine string
		var n int

		for {
			buff := make([]byte, bufferSize)
			if isCompressed {
				n, err = zr.Read(buff)
			} else {
				n, err = r.Read(buff)
			}

			if err != nil && err != io.EOF {
				if isCompressed {
					zr.Close()
				}
				return nil, nil, fmt.Errorf("got error during read %s: %s", file, err.Error())
			}

			if logFile == nil && (string(buff[:n]) == "done\n" || string(buff[:n]) == "exit\n") {
				break
			}

			leftLine += string(buff[:n])
			lines := strings.SplitAfter(leftLine, "\n")

			for _, line := range lines {
				if strings.HasSuffix(line, "\n") || err == io.EOF {
					if len(line) > 0 {
						if err == io.EOF && !strings.HasSuffix(line, "\n") {
							line += "\n"
						}

						ltsv, _ := parseLtsv(line)

						if baseKey != "" {
							resultKey := ltsv[baseKey]
							resultValue := ltsv[targetKey]

							if resultKey == "" {
								continue
							}

							// update maxKeyLength
							if len(resultKey) > maxKeyLength {
								maxKeyLength = len(resultKey)
							}

							// make key for new value and increase counter
							if _, exist := results[resultKey]; !exist {
								results[resultKey] = &result{
									sum:   0,
									avg:   0,
									count: 1,
								}
							} else {
								results[resultKey].count++
							}

							if targetKey != "" {
								r, err := strconv.Atoi(resultValue)
								if err == nil {
									// if targetKey is integer, set sum and avg for each base key values
									results[resultKey].sum += uint64(r)
									results[resultKey].avg = results[resultKey].sum / results[resultKey].count
								}
							}
						}
					}

					leftLine = ""
				} else {
					leftLine = line
				}
			}

			if err == io.EOF {
				break
			}
		}

		if isCompressed {
			zr.Close()
		}
		if logFile != nil {
			logFile.Close()
		}
	}

	for key, value := range results {
		switch sortKey {
		case "SUM":
			resultKeys = append(resultKeys, sortedKey{key, value.sum})
			break
		case "AVG":
			resultKeys = append(resultKeys, sortedKey{key, value.avg})
			break
		case "CNT", "COUNT":
			resultKeys = append(resultKeys, sortedKey{key, value.count})
			break
		}
	}

	sort.SliceStable(resultKeys, func(i, j int) bool {
		if reverseSort {
			return resultKeys[i].Value < resultKeys[j].Value
		}
		return resultKeys[i].Value > resultKeys[j].Value
	})

	return results, resultKeys, nil
}

func parseLtsv(line string) (map[string]string, []string) {
	columns := strings.Split(line, "\t")
	ltsv := make(map[string]string)
	var keys []string

	for _, column := range columns {
		lv := strings.SplitN(column, ":", 2)
		if len(lv) < 2 {
			continue
		}
		key, param := strings.TrimSpace(lv[0]), strings.TrimSpace(lv[1])
		ltsv[key] = param
		keys = append(keys, key)
	}

	return ltsv, keys
}
