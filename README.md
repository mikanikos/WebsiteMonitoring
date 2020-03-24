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
The app could be improved in several aspects. First of all, there are not enough tests (for lack of time I couldn't include more) but several manual tests have been carried out to make sure that the basic functionalities work as expected. 

In second place, the application can be structured in a better way in order to be even more modular. The current organization is severely split in multiple packages in order to benefit of different indipendent components and make them interact coherently. However, some functionalities could be better grouped and simplified (for example, websiteManager and statsManager are strictly correlated and there might be a better way to improve the design of a general stats and website handler and processor).

The current organization allows to add more stats but not more alerts: alerts should include an interface that can give the possibility to extend and enrich the application with more user-defined alerts. Moreover, the current alert handler does a similar job to what a WebsiteStatsHandler does: a better option would be to make the processing of measures unique for both stats handlers and alerts so that there's no more the need to handle data similarly in different places.

This implementation processes data from the website monitor in FIFO order and assumes this order even in the other processing methods: this is not a real limitation but this assumption can be better defined in code and made more clear for future extensions. 

Stats are computed with aggregators of metrics which should result in efficient processing of measures but aggregators can be made even more general in order to allow the creation of more metrics combinations.

User input is not checked and this can be improved by adding some security checks when parsing the config file.

More stats could be added in order to better monitor websites' performance.

An local database like SQLite could be used to make the data persistent and save the aggregates: this would make the application more fault-tolerant.

The monitor puts a timeout for the response of a website which is half of the check interval. Initially, this was a user parameter but, since the timeout must be strictly lower than the check interval and there were no checks on input data, I decided to remove it. A better way could be found to intergate this addition and make it user-configurable in a safe way.  

From a performance point of view, there might be issues if the tool is scaled to more websites and stats. Currently, a new thread handles website monitoring and another thread handles measure processing. Moreover, a new thread is spawned for each stats that is used in the application. This choice was made to have a more flexible tool with separate indipendent components. However, a better design could make the application scale better even though no benchmarks have been made to assess the current performance.    

Another interesting addition would be making the user configuration parameters dynamics, i.e. the user can dynamically change these parameters while running the program. The current design would probably allow this improvement with few modifications.

The current user interface is limited but does the job. A first improvement would be having a dedicated webpage that can be better suit user needs and allows more flexibility in terms of paramaters choice and user experience.
