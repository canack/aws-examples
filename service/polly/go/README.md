# TTS with using AWS Polly

### Build:

```shell
git clone https://github.com/canack/aws-examples.git
cd aws-examples/service/polly/go
go build
```

### Usage:
```shell
./polly-example "Hello, I am Polly. I am a text-to-speech service."

# execute ls command after the execution of the program.
# There should be a file named "speech.mp3" in the current directory.
```

### Feel to free to contribute. :)