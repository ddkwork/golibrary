package pather

type (
	Interface interface {
		//Fn() (ok bool)
	}
	Input struct {
		LowerFileName string
		UpperFileName string
		EmbedPrefix   string
		ObjectSuffix  string
		Path          string
	}
	OutPut struct {
		Embed  string
		Kind   string
		Object string
		Method string
	}
	Object struct {
		Input
		OutPut
	}
)

func New(input Input) *Object {
	return &Object{
		Input: input,
		OutPut: OutPut{
			Embed:  input.EmbedPrefix + input.UpperFileName,
			Kind:   input.LowerFileName + "Kind",
			Object: input.LowerFileName + input.ObjectSuffix,
			Method: input.UpperFileName + input.ObjectSuffix,
		},
	}
}

//			in.mapEmbedPath["bmp_"+lower] = path
//			in.mapKindPath[lower+"Kind"] = path
//			in.mapObjectPath[lower+"Button "] = path
//			in.mapMethodPath[upper+"Button "] = path
