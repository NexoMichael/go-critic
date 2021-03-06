package checker_test

func closeFiles(files []*string) {
	f := func(s *string) {}

	/// key variable occurs more then once in the loop; consider using for _, value := range files
	for i := range files {
		if files[i] != nil {
			f(files[i])
		}
	}
}

func sliceLoop(files []*string) {
	f := func(s *string) {}

	/// key variable occurs more then once in the loop; consider using for _, value := range files
	for k := range files[:] {
		if files[k] != nil {
			f(files[k])
		}
	}
}
