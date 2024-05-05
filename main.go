package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    data, lenght := parseArguments(os.Args[1:])
    firstGen := toBoolSlice(data)
    printOutOctal(ruleThirty(firstGen, lenght))
}

func parseArguments(args []string) (string, int) {

    data := ""
    lenght := 128
    exist := [...]bool{false, false, false}
    helpMessage := fmt.Sprintf("crycellaton is a Cryptographic Cellular Automaton based on rule 30.\nUsage:\n\tcrycellaton [OPTION] argument\nOptions:\n\t'-l', '--lenght'\tchange the lenght of the final output. Default is %d.\n\t'-f', '--file'\treads file instead of argment.", lenght)

    for i := 0 ; i < len(args) ; i = i + 1 {
        
        if i == len(args) {
            break
        }

        switch args[i] {
            case "-l", "--lenght":
                if i + 1 == len(args) {
                    abort(fmt.Sprintf("< ! >\tThe '-l' [OPTION] require a decimal value, got nothing.\n%s", helpMessage))
                } else if exist[1] == true {
                    abort(fmt.Sprintf("< ! >\tThe '-l' [OPTION] appear twice.\n%s", helpMessage))
                } else {
                    lenght, _ = strconv.Atoi(args[i + 1])
                    i = i + 1
                    exist[1] = true
                }

            case "-f", "--file":
                if i + 1 == len(args) {
                    abort(fmt.Sprintf("< ! >\tThe '-f' [OPTION] require a file path, got nothing.\n%s", helpMessage))
                } else if exist[0] == true {
                    abort(fmt.Sprintf("< ! >\tThe '-f' [OPTION] isn't compatible with [DEFAULT] argument.\n%s", helpMessage))
                } else if exist[2] == true {
                    abort(fmt.Sprintf("< ! >\tThe '-f' [OPTION] appear twice.\n%s", helpMessage))
                } else {
                    // we just read the inner content of a given file, maybe it's not ideal ?
                    data = fileReader(args[i + 1])
                    i = i + 1
                    exist[2] = true
                }

            case "-h", "--help":
                abort(helpMessage)

            default:
                if exist[0] == true {
                    abort(fmt.Sprintf("< ! >\tThe [DEFAULT] argument appear twice. Or maybe you misstyped an [OPTION] ?\n%s", helpMessage))
                } else if exist[2] == true {
                    abort(fmt.Sprintf("< ! >\tThe [DEFAULT] argument isn't compatible with the '-f' [OPTION].\n%s", helpMessage))
                } else {
                    data = args[i]
                    exist[0] = true
                }
        }
    }

    if exist[0] == false && exist[2] == false {
        abort(helpMessage)
    }

    return data, lenght
}

func fileReader(path string) string {

    data, err := os.ReadFile(path)

    if err != nil {
        // no need to close, we are not reading the file, so we just throw error an quit the whole program, gc will do the rest
        abortError(err)
    }

    return string(data)
}

func ruleThirty(input []bool, iterations int) []int {

    defaultRuleSet := [8]bool{false, false, false, true, true, true, true, false}
    gen := make([]bool, 0)
    seed := make([]bool, 0)
    result := make([]bool, 0)

    // originally this was iterations * 2, to prevent bias, it was changed to make 128 iterations before setting the final bits
    // there is a bias with this method aswell, each generated sequence, no matter of the lenght specified by the user will generate the same result. generating a smaller result will be only troncate the final one and it's kinda bad
    // ont thing to consider is generating without keeping the bits by generation but with another method that i don't know yet
    // here i tried to remove the bias by generating a seed with 128 generations at minimum
    // then i 
    for i := 0 ; i <= 128 + iterations ; i = i + 1 {
        if i == 0 {
            gen = nextGeneration(defaultRuleSet[:], input)
        } else {
            gen = nextGeneration(defaultRuleSet[:], gen)
        }
        // the bias is done here
        if i > 128 {
            seed = append(seed, modTwo(intSliceToInt(boolSliceToIntSlice(gen))))
        }
    }
    // it seems that re-generating from a seed wich was biased is not worth
    for i := 0 ; i <= 128 ; i = i + 1 {
        if i == 0 {
            result = nextGeneration(defaultRuleSet[:], seed)
        } else {
            result = nextGeneration(defaultRuleSet[:], result)
        }
    }

    return boolSliceToIntSlice(result)
}

func toBoolSlice(s string) []bool {

    tempStr := ""
    ret := make([]bool, 0)

    for _, c := range s {
        tempStr = fmt.Sprintf("%s%b", tempStr, int(rune(c)))
    }

    for _, c := range tempStr {
        tempBool, _ := strconv.ParseBool(string(c))
        ret = append(ret, tempBool)
    }

    return ret
}

func nextGeneration(ruleset []bool, world []bool) []bool {

    if len(world) <= 3 {
        abort(fmt.Sprintf("the next generation can't be processed ! world lengh is too low, got [%d]", len(world)))
    }

    ret := make([]bool, 0)

    for index, cell := range world {

        switch index {
            case 0:
                neighborhood := [3]bool{ world[len(world) - 1], cell, world[index + 1] }
                ret = append(ret, ruleCheck(ruleset, neighborhood))
            case len(world) - 1:
                neighborhood := [3]bool{ world[index - 1], cell, world[0] }
                ret = append(ret, ruleCheck(ruleset, neighborhood))
            default:
                neighborhood := [3]bool{ world[index - 1], cell, world[index + 1] }
                ret = append(ret, ruleCheck(ruleset, neighborhood))
        }

    }

    return ret
}

func ruleCheck(ruleset []bool, neighborhood [3]bool) bool {

    rules := [8][3]bool{
        {true,  true,  true },
        {true,  true,  false},
        {true,  false, true },
        {true,  false, false},
        {false, true,  true },
        {false, true,  false},
        {false, false, true },
        {false, false, false},
    }

    switch neighborhood {
        case rules[0]:
            return ruleset[0]
        case rules[1]:
            return ruleset[1]
        case rules[2]:
            return ruleset[2]
        case rules[3]:
            return ruleset[3]
        case rules[4]:
            return ruleset[4]
        case rules[5]:
            return ruleset[5]
        case rules[6]:
            return ruleset[6]
        case rules[7]:
            return ruleset[7]
        default:
            abort(fmt.Sprintf("can't verify neighborhood with ruleset ! rules are [%v], got neighborhood of [%v]", rules, neighborhood))
    }

    return false
}

func boolSliceToIntSlice(bs []bool) []int {

    ret := make([]int, 0)

    for _, b := range bs {
        if b == true {
            ret = append(ret, 1)
        } else {
            ret = append(ret, 0)
        }
    }

    return ret
}

func modTwo(v int) bool {
    /*
    if v % 2 == -1 {
        return 1
    } else {
        return v % 2
    }
    */
    if v % 2 == -1 || v % 2 == 1 {
        return true
    } else {
        return false
    }

}

func intSliceToInt(sl []int) int {

    t := ""
    ret := 0

    for _, v := range sl {
        t = t + strconv.Itoa(v)
        temp, _ := strconv.Atoi(t)
        ret = ret + temp
    }

    return ret
}

func intSliceToSingleInt(sl []int) int {

    temp := ""

    for _, v := range sl {
        temp = temp + strconv.Itoa(v)
    }

    ret, _ := strconv.Atoi(temp)

    return ret
}

func printOutBinary(row []int) {

    for i, v := range row {
        if i == len(row) - 1 {
            fmt.Printf("%d\n", v)
            break
        }
        fmt.Print(v)
    }

}

func printOutOctal(row []int) {

    i := 0
    tempSlice := make([][]int, 0)

    for ; i <  len(row) - 1 ; i = i + 7 {
        if i + 7 > len(row) - 1 {

            for j := 0 ; j == 6 ; j = j + 1 {
                if i + j > len(row) - 1 {
                    continue
                } else {
                    tempSlice = append(tempSlice, row[i:i+j])
                }
            }

        } else {
            tempSlice = append(tempSlice, row[i:i+7])
        }
    }

    ret := ""

    for _, sl := range tempSlice {
        tempString := ""

        for _, v := range sl {
            tempString = fmt.Sprintf("%s%s", tempString, strconv.Itoa(v))
        }

        tempInt, _ := strconv.Atoi(tempString)
        ret = fmt.Sprintf("%s%x", ret, tempInt)
    }

    fmt.Println(ret)
}

func abort(s string) {
    fmt.Println(s)
    os.Exit(1)
}

func abortError(err error){
    fmt.Println(err)
    os.Exit(1)
}
