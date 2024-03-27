package stream

// GenI18n 在项目初始化的地方使用init,gcs 是需要在model\gurps的位置初始化，个人项目应该是在main函数设置即可
//
//	func init() {
//		i18n.Dir = "."
//	}
//
// todo 调用nlp翻译引擎自动翻译生成文件的value
// k:"Note"
// v:"备注"
// 以适应项目更新后每次手动翻译的繁琐问题，但是要解决翻译引擎不准确的问题
//
//	func transLateToChinese() {
//		const name = "zh-CN.i18n"
//		lines, ok := stream.ReadFileToLines(name)
//		if !ok {
//			return
//		}
//		t := translate.New()
//		for i, line := range lines {
//			if strings.Contains(line, "v:") {
//				split := strings.Split(line, ":")
//				v := split[1]
//				//v:"Toggle Equipped"
//				//todo call nlp
//				comment := t.Translate(v)
//				println(i)
//				println(v)
//				println(comment)
//				if !strings.Contains(comment, `"`) {
//					comment = strconv.Quote(comment)
//				}
//				lines[i] = "v:" + comment
//				if i == 40 {
//					//break
//				}
//			}
//		}
//		linesToString := stream.New("").LinesToString(lines)
//		stream.WriteTruncate(name, linesToString)
//	}
//
//	func GenI18nWithTransLateToChinese(path string) {
//		GenI18n(path)
//		//transLateToChinese()
//	}
//func GenI18n(path string) {
//	//cmdline.CopyrightStartYear = "2016"
//	//cmdline.CopyrightHolder = "Richard A. Wilkes"
//	//cmdline.License = "Mozilla Public License 2.0"
//	//cl := cmdline.New(true)
//	//cl.UsageSuffix = "<path> [path...]"
//	//cl.Description = i18n.Text("Generates a template for a localization file from source code.")
//
//	outPath := i18n.Locale() + ".i18n"
//	//cl.NewGeneralOption(&outPath).SetSingle('o').SetName("output").SetArg("path").SetUsage("The output file")
//	//args := cl.Parse(os.Args[1:])
//	args := []string{path}
//	//if outPath == "" {
//	//	cl.FatalMsg(i18n.Text("The output file may not be an empty path."))
//	//}
//	//if len(args) == 0 {
//	//	cl.FatalMsg(i18n.Text("At least one path must be specified."))
//	//}
//	kv := make(map[string]string)
//	fileSet := token.NewFileSet()
//	for _, pathArg := range args {
//		var err error
//		if pathArg, err = filepath.Abs(pathArg); err == nil {
//			walkErr := filepath.WalkDir(pathArg, func(path string, fi os.DirEntry, err error) error {
//				if err != nil {
//					return err
//				}
//				if !fi.IsDir() && filepath.Ext(path) == ".go" {
//					fmt.Println(path)
//					var file *ast.File
//					if file, err = parser.ParseFile(fileSet, path, nil, 0); err != nil {
//						fmt.Fprintln(os.Stderr, err)
//						if !mylog.Error(err) {
//							return err
//						}
//					}
//					const (
//						LookForPackageState = iota
//						LookForTextCallState
//						LookForParameterState
//					)
//					state := LookForPackageState
//					ast.Inspect(file, func(node ast.Node) bool {
//						switch x := node.(type) {
//						case *ast.Ident:
//							switch state {
//							case LookForPackageState:
//								if x.Name == "i18n" {
//									state = LookForTextCallState
//								}
//							case LookForTextCallState:
//								if x.Name == "Text" {
//									state = LookForParameterState
//								} else {
//									state = LookForPackageState
//								}
//							default:
//								state = LookForPackageState
//							}
//						case *ast.BasicLit:
//							if state == LookForParameterState {
//								if x.Kind == token.STRING {
//									var v string
//									if v, err = strconv.Unquote(x.Value); err != nil {
//										fmt.Fprintln(os.Stderr, err)
//									} else {
//										kv[v] = v
//									}
//								}
//							}
//							state = LookForPackageState
//						case nil:
//						default:
//							state = LookForPackageState
//						}
//						return true
//					})
//				}
//				return nil
//			})
//			if walkErr != nil {
//				fmt.Fprintln(os.Stderr, walkErr)
//			}
//		} else {
//			fmt.Fprintln(os.Stderr, err)
//		}
//	}
//
//	keys := make([]string, 0, len(kv))
//	for key := range kv {
//		keys = append(keys, key)
//	}
//	sort.Slice(keys, func(i, j int) bool {
//		return txt.NaturalLess(keys[i], keys[j], true)
//	})
//	out, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
//	if !mylog.Error(err) {
//		fmt.Fprintf(os.Stderr, "Unable to create '%s'.\n", outPath)
//		return
//	}
//	//# 翻译替换所有v:"xxx"为v:"翻译的中文内容"，保留原始换行格式和k:"xxx"，发成代码块过来
//	fmt.Fprintf(out, `
//# Generated on %v
//#
//# Key-value pairs are defined as one or more lines prefixed with "k:" for the
//# key, followed by one or more lines prefixed with "v:" for the value. These
//# prefixes are then followed by a quoted string, using escaping rules for Go
//# strings where needed. When two or more lines are present in a row, they will
//# be concatenated together with an intervening \n character.
//#
//# Do NOT modify the 'k' values. They are the values as seen in the code.
//#
//# Replace the 'v' values with the appropriate translation.
//`, time.Now().String(time.RFC1123))
//	for _, key := range keys {
//		fmt.Fprintln(out)
//		for _, p := range strings.Split(key, "\n") {
//			if _, err = fmt.Fprintf(out, "k:%q\n", p); err != nil {
//				fmt.Fprintln(os.Stderr, err)
//				if !mylog.Error(err) {
//					return
//				}
//			}
//		}
//		for _, p := range strings.Split(key, "\n") {
//			if _, err = fmt.Fprintf(out, "v:%q\n", p); err != nil {
//				fmt.Fprintln(os.Stderr, err)
//				if !mylog.Error(err) {
//					return
//				}
//			}
//		}
//	}
//	if err = out.Close(); err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		if !mylog.Error(err) {
//			return
//		}
//	}
//}
