package outputGenerators

type FileGenerator interface {
	Generate(plugins map[string]string) []byte
}

func SetOutputGenerator(fg FileGenerator) FileGenerator {
	return fg
}
