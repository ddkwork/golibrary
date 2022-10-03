package myc2go_test

//
//func TestCpp2Go(t *testing.T) {
//	o := myc2go.NewObj()
//	o.Src("D:\\codespace\\gui\\sdk\\HyperDbgDev\\hyperdbg").
//		Dst("binding").
//		ExpandPath("miscellaneous\\constants", ".txt").
//		Back().
//		Convert()
//	//o.Format()
//}
//
//func TestUnsafe(t *testing.T) {
//	sizeofUINT32 := unsafe.Sizeof(uint32(0))
//	mylog.Assert(t).Equal(uint32(sizeofUINT32), uint32(4))
//}
//func Test3(t *testing.T) {
//	dir := filepath.Dir("back\\HyperDbgDev\\hyperdbg\\include\\SDK\\Headers\\Constants.h.back")
//	pkgName := filepath.Base(dir)
//	println(pkgName)
//}
//
//func Test2(t *testing.T) {
//	p := "../HyperDbgDev/hyperdbg"
//	slash := filepath.ToSlash(p)
//	_, after, found := strings.Cut(slash, `/`)
//	if !found {
//		panic("!found")
//	}
//	join := filepath.Join("back", after)
//	println(join)
//}
