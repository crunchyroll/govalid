    func validateBool(fldname string) {
    	fmt.Printf("\tret.%s, err = strconv.ParseBool(data[\"%s\"])\n", fldname, fldname)
    	fmt.Printf("\tif err != nil {\n")
    	fmt.Printf("\t\treturn nil, err\n")
    	fmt.Printf("\t}\n")
    }

    case "bool":
    	validateBool(nam)
    	needsStrconv = true
