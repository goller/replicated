package cmd

import (
	"github.com/pkg/errors"
	"github.com/replicatedhq/replicated/cli/print"
	"github.com/spf13/cobra"
)

func (r *runners) InitCustomersLSCommand(parent *cobra.Command) *cobra.Command {
	customersCmd := &cobra.Command{
		Use:   "ls",
		Short: "list customers",
		Long:  `list customers`,
		RunE: r.listCustomers,
	}
	parent.AddCommand(customersCmd)

	return customersCmd
}

func (r *runners) listCustomers(_ *cobra.Command, _ []string) error {

	customers, err := r.api.ListCustomers(r.appID, r.appType)
	if err != nil {
		return errors.Wrap(err, "list customers")
	}

	return print.Customers(r.stdoutWriter, customers)
}
