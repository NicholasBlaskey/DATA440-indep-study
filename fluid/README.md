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

10 samples of images were generated for each combination.

b2 = framesPerSample = 275, DT = 0.01, saveEvery = 5, width = 128

45.794378837s, 3.8M

b3 = framesPerSample = 275, DT = 0.01, saveEvery = 5, width = 256

46.004321232s, 9.2M

b1 = framesPerSample = 275, DT = 0.01, saveEvery = 5, width = 512

55.79745923s, 23M

b4 = framesPerSample = 275, DT = 0.01, saveEvery = 5, width = 1024

2m41.497036213s, 70M

The time is relatively constant across the actual rendering because of the way the solver uses framebuffer. The real increase in time is the making of PNG files and getting data out of the framebuffer for the GPU encoding.

Let's just assume we are gonna generate 10,000 samples.

128 would need 380MB
256 would need 920MB
512 would need 2.3GB
1024 would need 7GB

512 gives really good detail. 256 also gives a peak of the detail but not as great.

By converting the encoding to the PNG with the same settings as b4 we got a time of 45.836481204s. Now we are limited by the GPU instead of CPU bound likely.
