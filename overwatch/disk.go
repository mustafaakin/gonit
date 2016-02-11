package overwatch

import (
	"os"
	"log"
)

// PrepareDisk prepares the disk for the required folders 
func PrepareDisk(){
	folders := []string{"aq", "bq"}
	

	

	
	for _, folder := range folders {
		// A sensible way to create folders
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			log.Println("An error occured while preparing folder:", folder, err)
		}
	} 
}

// MountCGroupfs mounts the cgroup filesystem
func MountCGroupfs(){
	// TODO: Do it properly with syscalls
}