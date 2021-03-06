package main

import (
	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/vmware/govmomi/object"
	"fmt"
	"log"
	"time"
)

type StepShutdown struct {
}

func (s *StepShutdown) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	d := state.Get("driver").(*Driver)
	vm := state.Get("vm").(*object.VirtualMachine)

	ui.Say("Shut down VM...")

	err := d.StartShutdown(vm)
	if err != nil {
		state.Put("error", fmt.Errorf("Cannot shut down VM: %v", err))
		return multistep.ActionHalt
		}

	timeoutvalue := time.Second * 300
	log.Printf("Waiting max %s for shutdown to complete", timeoutvalue)
	err = d.WaitForShutdown(vm, timeoutvalue)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say("VM stopped")
	return multistep.ActionContinue
}

func (s *StepShutdown) Cleanup(state multistep.StateBag) {}
