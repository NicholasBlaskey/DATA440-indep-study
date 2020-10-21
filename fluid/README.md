### Running

```
go run main.go
```

### Benchmarks

There are a couple of parameters to tune in image generation.

#### Width and Height

We are gonna ensure width = height. This is just the output dimensions of the image.

#### dt

This is the time step used in the fluid algorithm. Note we don't want this timestep to be too large because the algorithm relies on it being sufficently small. This is the time between each frame or update iteration.

#### saveEvery

This is an integer that says we will save every nth frame. This allows us to get over the dt being too large.

#### framesPerSample

This is how many frames we are going to use in sample. This we want to be large enough so that the fluid always goes away for the next frame but small enough so we don't waste time recording blank images. This should essentially be a function of dt.

### Results

#### First let's look at the effect of setting width and height. 

10 samples of images were generated for each combination. This was done before further optimizations later on.

| framesPerSample | DT   | saveEvery | Width | time (seconds) | size (MB) | size for 10,000 samples (GB) |
| --------------- | ---- | ----------| ----- | -------------- | --------- | ---------------------------- |
| 275             | 0.01 | 5         | 128   | 45.79          | 3.8       | 0.380                        |
| 275             | 0.01 | 5         | 256   | 46.00          | 9.2       | 0.920                        |
| 275             | 0.01 | 5         | 512   | 55.80          | 23        | 2.3                          |
| 275             | 0.01 | 5         | 1024  | 161.50         | 70        | 7.0                          |


The time is relatively constant across the actual rendering because of the way the solver uses framebuffer. The real increase in time is the making of PNG files and getting data out of the framebuffer for the GPU encoding.

512 gives really good detail. 256 also gives a hint of the detail but not as great.

#### Multithreaded and GPU

By converting the encoding to the PNG with the same settings as b4 we got a time of 45.836481204s. Now we are limited by the GPU instead of CPU bound likely. By running on my other laptop which has a GPU I got 36.4477372s.


#### Paramters used

We are going to use 512px by 512px images because of the detail that is allowed to be captured and we can scale them down later.

We are going to stick with a dt of 0.01 because it gives a reasonable result.

We are going to stick to saving every 3rd frame because it is the lowest amount of frames saved that can be handled by my computers CPUs without slowing down the GPU bound part. It also gives good enough results.

My laptop takes 24.6483853s to run 10 samples.