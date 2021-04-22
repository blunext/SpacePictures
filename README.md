Service gets pictures URLs from NASA API Picture of the Day.

Sample request:
```
http://localhost:8080/pictures?start_date=2020-10-20&end_date=2020-10-22
```

Sample response
```
{“urls”: ["https://apod.nasa.gov/apod/image/2008/AlienThrone_Zajac_3807.jpg", ...]}
```


To run it
```
Tu go build
./SpacePictures 
```

or run it via Docker
```
./run.sh
```


