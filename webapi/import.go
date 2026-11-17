package webapi

import "github.com/grafana/sobek"

func installImport(r *Realm) {
	r.VM.SetGetImportMetaProperties(func(mr sobek.ModuleRecord) []sobek.MetaProperty {
		return []sobek.MetaProperty{
			{
				Key: "url",
				Value: r.VM.ToValue(
					r.Doc.GetModuleReferrerUrl(mr).String(),
				),
			},
			{
				Key: "resolve",
				Value: r.VM.ToValue(func(call sobek.FunctionCall) sobek.Value {
					specifier := call.Argument(0).String()
					url, err := r.Doc.ResolveModuleSpecifier(mr, specifier)
					if err != nil {
						panic(r.VM.NewTypeError("failed to resolve module specifier %v: %v", specifier, err))
					}

					return r.VM.ToValue(url.String())
				}),
			},
		}
	})

	// TODO: fetch modules asynchronously
	r.VM.SetImportModuleDynamically(func(referrer interface{}, specifier sobek.Value, promiseCapability interface{}) {
		mod, err := r.Doc.ResolveImportedModule(referrer, specifier.String())
		if err != nil {
			r.VM.FinishLoadingImportModule(referrer, specifier, promiseCapability, nil,
				r.VM.NewTypeError("failed to resolve module specifier %v: %v", specifier, err),
			)
			return
		}

		r.VM.FinishLoadingImportModule(referrer, specifier, promiseCapability, mod, nil)
	})
}
