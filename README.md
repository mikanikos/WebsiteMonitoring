# WebsiteMonitoringTool
Website availability &amp; performance monitoring

## Run the application

To download the project, simply use the following command 

```
go get -v github.com/mikanikos/WebsiteMonitoring
```

Once you downloaded the repository, go to the main directory (where the main.go is located), compile the project and run it:

```
go build
go run main.go -config="config.yaml"
```

Please note that the config file is located in the configs folder and you don't need to specify the path: by default all the configurations files must be located in the configs folder.

## Improvements
The app could be improved on several aspects: first of all, there are not enough tests (for lack of time I couldn't include more); in second place, the app can be structured in a better way in order to be even more modular. In fact, alerts should have an interface that could allow to have more alerts with the same interface and gives the possibility to extend the application with more user configuration options.

Stats are computed with aggregators which should be quite efficient but aggregators can be made more general in order to allow the creation of more combinations of other metrics.

User input is not checked and this can be improved by adding some checks when parsing the config file.
