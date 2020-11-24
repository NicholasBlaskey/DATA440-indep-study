### Install

I had trouble installing by doing the reccomended docs of just doing

```
pip install tensorflowjs
```

I needed to install miniconda then

```
conda create -n tfjs python=3.6.8
conda activate tfjs
pip install tensorflowjs
```

https://stackoverflow.com/questions/48720833/could-not-find-a-version-that-satisfies-the-requirement-tensorflow

https://stackoverflow.com/questions/50317081/tensorflowjs-converter-command-not-found

Then the JSON needed to be editing according to this post due to another error saying the layer functional was not supported

https://stackoverflow.com/questions/63143849/tensorflow-js-error-unknown-layer-functional

### Running

First convert the tensorflow generator.h5 pretrained model into the Layers format

```
tensorflowjs_converter keras --input_format generator.h5 ./layersGenerator/
```

Then run a simple server by doing

```
python3 -m http.server
```

Finally go to
```
127.0.0.1:8000
```

