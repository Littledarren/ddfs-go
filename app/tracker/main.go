package main

func main() {
	s := NewTrackerServiceImpl()
	s.ListenAndServe()
}
