run: install
	go run *.go

install: clean
	gcc -E HelloWorldData.c -o h.i -I ../../library/include
	gcc -S h.i -o h.s
	gcc -c h.s -o HelloWorldData.o

clean:
	rm -f h.i h.s HelloWorldData.o