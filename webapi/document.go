package webapi

func installDocument(r *Realm) {
	doc := r.VM.NewObject()

	r.VM.Set("document", doc)
}
