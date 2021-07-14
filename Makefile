FIVEAUDIO=./Five.mp3

Five: main.go FivePre.dmp3 FiveLoop.dmp3
	go build -o ./Five main.go
	$(MAKE) clean

FivePre.dmp3 FiveLoop.dmp3: transport.go $(FIVEAUDIO)
	go run transport.go $(FIVEAUDIO)

clean:
	$(RM) FivePre.dmp3 FiveLoop.dmp3
