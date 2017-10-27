package main

import "hrentabd"

type Storage interface {

	Save(tab hrentabd.HrenTab)
	Load() hrentabd.HrenTab

}
