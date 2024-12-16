# OpenSatMap data acquisition

## Description
This code is used to get images from google maps <b>given a GPS region (lat1, lon1, lat2, lon2) or a center GPS point and a Zoom level.</b>

The code is written in Go and uses the Google Maps Static API to get the images.

This code is used in [OpenSatMap (NeurIPS 2024 D&B Track)](https://github.com/OpenSatMap/OpenSatMap-offical) to collect the data. 
## Usage
Here are the steps to get the images:
1. Install Go. The Go version we used is 1.19.
``` bash
wget  https://go.dev/dl/go1.19.linux-amd64.tar.gz

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin
```
2. Download some dependencies.
``` bash
GO111MODULE=off go get github.com/ironsublimate/gomapinfer/common

GO111MODULE=off go get github.com/ironsublimate/gomapinfer/googlemaps
```
3. Fix the bugs in the code.
You should go to path/to/src/github.com/ironsublimate/gomapinfer/googlemaps foloder, modify coords.go, sat.go.
Use the following code to fix the bugs. (The case of ironsublimate is changed from upper to lower case.)
``` go
import (
        "github.com/ironsublimate/gomapinfer/common"
)
```


4. Run the code.
``` bash
mkdir -p ./data/imagery
# get images from a GPS region
GO111MODULE=off go run1_sat.go KEY ./data/imagery
# get images from a GPS center and a zoom level
GO111MODULE=off go run1_sat_center.go KEY ./data/imagery
```

## Acknowledgement
The code is based on the code from [roadtracer](https://github.com/mitroadmaps/roadtracer/tree/master/dataset). 
Thanks to the authors for sharing the code.

## Citation
```
@article{zhao2024opensatmap,
  title={OpenSatMap: A Fine-grained High-resolution Satellite Dataset for Large-scale Map Construction},
  author={Zhao, Hongbo and Fan, Lue and Chen, Yuntao and Wang, Haochen and Jin, Xiaojuan and Zhang, Yixin and Meng, Gaofeng and Zhang, Zhaoxiang},
  journal={arXiv preprint arXiv:2410.23278},
  year={2024}
}
```
