package main

import (
	"fmt"
	"os"
	"strings"
)

var tFiles []string
var baseKey string
var baseParam string
var targetKey string
var operation int
var sortKey string
var reverseSort bool
var maxKeyLength int

const (
	summary = 1
	avarage = 2
	counter = 4
)

func init() {
	if len(os.Args) < 2 {
		Usage()
		os.Exit(1)
	}

	args := os.Args
	maxArgs := len(args)

	// for debug
	// baseKey = "host"
	// targetKey = "bytes_sent"
	// tFiles = append(tFiles, "logs/20200627/access.log_14")
	baseKey = ""
	targetKey = ""
	operation = 0
	sortKey = "CNT"
	reverseSort = false
	maxKeyLength = 0

	for i := 1; i < maxArgs; i++ {
		argv := args[i]

		switch argv {
		case "--base", "-B":
			if i+1 < maxArgs {
				baseKey = args[i+1]
				i++
			} else {
				Usage()
				os.Exit(1)
			}
			break
		case "--target", "-T":
			if i+1 < maxArgs {
				targetKey = args[i+1]
				i++
			} else {
				Usage()
				os.Exit(1)
			}
			break
		case "--operation", "-O":
			if i+1 < maxArgs {
				ops := strings.Split(args[i+1], ",")

				opMap := make(map[string]bool, 3)
				for _, op := range ops {
					op = strings.ToUpper(op)
					// handle duplicate operations
					if _, exist := opMap[op]; !exist {
						switch op {
						case "SUM":
							operation += summary
							break
						case "AVG":
							operation += avarage
							break
						case "CNT", "COUNT":
							operation += counter
							break
						case "ALL":
							operation = summary + avarage + counter
							break
						default:
							os.Stderr.WriteString("Not support operation! Only SUM, AVG, CNT can available.")
							Usage()
							os.Exit(1)
						}
					}
				}

				i++
			} else {
				Usage()
				os.Exit(1)
			}
			break
		case "--sort", "-S":
			if i+1 < maxArgs {
				arg := strings.ToUpper(args[i+1])
				switch arg {
				case "SUM", "AVG", "CNT", "COUNT":
					sortKey = arg
					i++
					break
				default:
					os.Stderr.WriteString("Not support sort key! Only SUM, AVG, CNT can available.")
					Usage()
					os.Exit(1)
				}
			} else {
				Usage()
				os.Exit(1)
			}
			break
		case "--rev", "-r", "-R":
			reverseSort = true
			break
		case "--help", "-h", "-H":
			Usage()
			os.Exit(1)
		default:
			tFiles = append(tFiles, argv)
		}
	}
}

// Usage show simple manual
func Usage() {
	fmt.Println("How to use web-proxy log analyzer")
	fmt.Println("ltsv-analyzer [OPTIONS] {PATH1} {PATH2} ... (default path : ./logs)")
	fmt.Println()
	fmt.Println("if [PATH] not defined, use stdin.")
	fmt.Println()
	fmt.Println("- OPTIONS")
	fmt.Println("    --base[-B]      : set base key")
	fmt.Println("    --target[-T]    : set target key")
	fmt.Println("    --operation[-O] : set operation (sum, avg, cnt[count], all)")
	fmt.Println("    --sort[-S]      : set sort key  (sum, avg, cnt[count])")
	fmt.Println("    --rev[-r|-R]    : reverse sort direction (default: DESCending)")
	fmt.Println()
	fmt.Println("  --help[-h|H]      : show usage")
}

func main() {
	results, resultKeys, err := ltsvAnalyze(tFiles)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	if len(resultKeys) != 0 {
		printResult(results, resultKeys)
	}

	os.Exit(0)
}

func printResult(results map[string]*result, resultKeys []sortedKey) {
	var banner string

	// set base key string by base key's length
	baseKeyString := fmt.Sprintf(fmt.Sprintf("%%-%ds", maxKeyLength), fmt.Sprintf(fmt.Sprintf("%%-%ds", maxKeyLength/2), fmt.Sprintf(fmt.Sprintf("%%%ds", maxKeyLength/2), baseKey)))

	if targetKey == "" {
		fmt.Printf("Print LOG COUNTER by BASE KEY [%s]\n", baseKey)
		operation = counter
	} else {
		fmt.Printf("Print results by BASE KEY [%s] and TARGET KEY [%s]\n", baseKey, targetKey)
	}

	// set banner by each operation types
	switch operation {
	case 1:
		banner = fmt.Sprintf("%s%20s", baseKeyString, fmt.Sprintf("SUM(%s)", targetKey))
		break
	case 2:
		banner = fmt.Sprintf("%s%20s", baseKeyString, fmt.Sprintf("AVG(%s)", targetKey))
		break
	case 3:
		banner = fmt.Sprintf("%s%20s%20s", baseKeyString, fmt.Sprintf("SUM(%s)", targetKey), fmt.Sprintf("AVG(%s)", targetKey))
		break
	case 5:
		banner = fmt.Sprintf("%s%20s%20s", baseKeyString, fmt.Sprintf("SUM(%s)", targetKey), "LOG COUNTER")
		break
	case 6:
		banner = fmt.Sprintf("%s%20s%20s", baseKeyString, fmt.Sprintf("AVG(%s)", targetKey), "LOG COUNTER")
		break
	case 7:
		banner = fmt.Sprintf("%s%20s%20s%20s", baseKeyString, fmt.Sprintf("SUM(%s)", targetKey), fmt.Sprintf("AVG(%s)", targetKey), "LOG COUNTER")
		break
	default:
		banner = fmt.Sprintf("%s%20s", baseKeyString, "LOG COUNTER")
		break
	}

	// print banner
	fmt.Println(banner)
	for i := 0; i < len(banner); i++ {
		fmt.Printf("-")
	}
	fmt.Println()

	// print results by each operation types
	for _, key := range resultKeys {
		switch operation {
		case 1:
			fmt.Printf(fmt.Sprintf("%%-%ds%s", maxKeyLength, "%20d\n"), key.Key, results[key.Key].sum)
			break
		case 2:
			fmt.Printf(fmt.Sprintf("%%-%ds%s", maxKeyLength, "%20d\n"), key.Key, results[key.Key].avg)
			break
		case 3:
			fmt.Printf(fmt.Sprintf("%%-%ds%s", maxKeyLength, "%20d%20d\n"), key.Key, results[key.Key].sum, results[key.Key].avg)
			break
		case 5:
			fmt.Printf(fmt.Sprintf("%%-%ds%s", maxKeyLength, "%20d%20d\n"), key.Key, results[key.Key].sum, results[key.Key].count)
			break
		case 6:
			fmt.Printf(fmt.Sprintf("%%-%ds%s", maxKeyLength, "%20d%20d\n"), key.Key, results[key.Key].avg, results[key.Key].count)
			break
		case 7:
			fmt.Printf(fmt.Sprintf("%%-%ds%s", maxKeyLength, "%20d%20d%20d\n"), key.Key, results[key.Key].sum, results[key.Key].avg, results[key.Key].count)
			break
		default:
			fmt.Printf(fmt.Sprintf("%%-%ds%s", maxKeyLength, "%20d\n"), key.Key, results[key.Key].count)
			break
		}
	}
}
