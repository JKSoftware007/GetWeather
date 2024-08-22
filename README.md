# GetWeather

This small example program will allow the user to submit coordinates to get the weather for the given location.

<h1>Install</h1>

git clone https://github.com/JKSoftware007/GetWeather.git

<h1>Build</h1>

cd into the cloned directory

    cd GetWeather

Build the code

    go build

The build will create GetWeather or GetWeather.exe depending whether you are on Linux or Windows.

<h1>Execute</h1>

    ./GetWeather

<h1>Test</h1>

The program runs on localhost at port 8080.  If on Windows, you may be prompted to allow access through Windows firewall.  On Linux you may have to disable your local firewall.

Once successfully running, you can use Postman on Windows or curl on Linux, or your favorite method for hitting the REST API.

<h3>Using curl</h3>

    curl --location 'http://127.0.0.1:8080/currentweather' \
    --header 'Content-Type: application/json' \
    --data '{
    "latitude" : 38.8894,
    "longitude" : -77.0352
    }'

<h3>Using Postman</h3>
Included in the root of the project is the exported Postman collection.  Import **GetWeatherExample.json** and use the GetWeather POST.
