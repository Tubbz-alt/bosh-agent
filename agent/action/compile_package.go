package action

import (
	"errors"

	boshmodels "github.com/cloudfoundry/bosh-agent/agent/applier/models"
	boshcomp "github.com/cloudfoundry/bosh-agent/agent/compiler"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type CompilePackageAction struct {
	compiler boshcomp.Compiler
}

func NewCompilePackage(compiler boshcomp.Compiler) (compilePackage CompilePackageAction) {
	compilePackage.compiler = compiler
	return
}

func (a CompilePackageAction) IsAsynchronous() bool {
	return true
}

func (a CompilePackageAction) IsPersistent() bool {
	return false
}

func (a CompilePackageAction) IsLoggable() bool {
	return true
}

func (a CompilePackageAction) Run(blobID, digestString, name, version string, deps boshcomp.Dependencies) (val map[string]interface{}, err error) {
	pkg := boshcomp.Package{
		BlobstoreID: blobID,
		Name:        name,
		Sha1:        digestString,
		Version:     version,
	}

	modelsDeps := []boshmodels.Package{}

	for _, dep := range deps {
		modelsDeps = append(modelsDeps, boshmodels.Package{
			Name:    dep.Name,
			Version: dep.Version,
			Source: boshmodels.Source{
				Sha1:        dep.Sha1,
				BlobstoreID: dep.BlobstoreID,
			},
		})
	}

	uploadedBlobID, uploadedDigest, err := a.compiler.Compile(pkg, modelsDeps)
	if err != nil {
		err = bosherr.WrapErrorf(err, "Compiling package %s", pkg.Name)
		return
	}

	result := map[string]string{
		"blobstore_id": uploadedBlobID,
		"sha1":         uploadedDigest.String(),
	}

	val = map[string]interface{}{
		"result": result,
	}
	return
}

func (a CompilePackageAction) Resume() (interface{}, error) {
	return nil, errors.New("not supported")
}

func (a CompilePackageAction) Cancel() error {
	return errors.New("not supported")
}
