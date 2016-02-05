package main

//go:generate yagi -tem=./tem/list.go -gen=int64;int32 -pkg=main
//go:generate yagi -tem=./tem/map.go -gen=string,*int64;string,*int32 -pkg=main
//go:generate yagi -tem=./mmap/mmap.go -gen=string,int64;string,int32 -pkg=main

func main() {

}
