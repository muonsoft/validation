# Validate locally before checking references

Implement Input.Validate and Service.Handle. Name, First and Second are required (ErrIsBlank at name/first/second). Collect all local errors without calling Repository.Exists or Save. If local validation succeeds, check both references, collecting ErrMissing at first and second when absent. On lookup failure preserve the error and stop. Save exactly once only after all validation succeeds. Preserve the supplied context and validator.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
