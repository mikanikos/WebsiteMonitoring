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

## Design overview
The goal of this implementation was not to make something that "just works" but to design the challenge in the best way from a performance and extendability point of view. Of course, it would have been easier to just collect measures, store all them in an array from the beginning and compute stats when neededl however, in my opinion this approach would not allow an efficient processing and storage of meaningful data and, regardless the simiplicity, it would have not scale to many more website and stats.

For this reason, I decide to make this tool as much configurable as possible by giving the user the possibility to set the website, stats and alert configuration from a single config file that must be located in the configs folder. It is possible to add as many websites as preferred and add as many periodic stats as preferred. Regarding alert, currently the application supports only the availability alert and it's possible to only modify the parameters of this one (so, multiple alerts are not supported). All parameters must be expressed in seconds in the config file.

The Monitor (monitor package) is an indipendent component that is resposible to monitor and collect measures from a single website periodically as specified by the user configuration. At every interval, it simply does a GET request and waits for a response for a maximum timeout. The timeout is set purposely to half of the duration of the check interval in order to avoid the case of waiting forever for a response and risking to create a bottleneck when processing incoming measures. 

New measures for a website are sent by the monitor to a channel that is processed by the WebsiteManager. A WebsiteManager is responsible to start the monitor, collect new measures from the monitor channel and coordinate the computations by interacting with other components. A WebsiteManager will use WebsiteStatsHandler instances to compute a stats for a specific stats and website configuration and will call the AvailabilityAlert methods for updating the alert as well.

An important design choice has been made regarding how to split the processing of measures. From a performance point of view, I decided to use aggregated metrics (aggregators) in order to compute stats that concern a specific time interval. This decision led me to split the processing of stats for the same website depending on the intervals chosen by the user. For example, if the user chose to have 2 stats (let's say one for the past 20 minutes and another one for the past 40 minutes), there will be metric aggregators for each stat because the time interval is different and the aggregated metric must take care of processing only measures after a certain time. This can be good for performance but it might not scale better as I thought, probably better options are present.

In order to process stats in an aggregated way, a WebsiteStatsHandler is used for each website and for each stats configuration specified by the user. A WebsiteStatsHandler will take care to process new measures at every update by purging outdated data (measures that are outside the time interval of the stats configuration) and updating the aggregators with new data. Aggregators metrics are defined in the metrics folder and take care of computing stats based on the data in input and output. For example, when I have an outdated measure I want to remove that measure from the metric so that is not taken anymore into consideration. On the other hand, when I have a brand-new measure, I'll simply add it to aggregator which will include in the computation of the metric.

Aggregators are a good way to avoid long computations when requested to display stats (think about very large interval for displaying stats and thousands of measures). This allows to keep track only of the data that is specifically needed for a specific website and a specific stats configuation and at the same time simplify the whole process.

To display stats, a StatsManager is used to coordinate the collection of metrics from the different WebsiteStatsHandlers and to call the ui. This component can be seen as a Presenter or Controller that will simply call methods from the logic of the application, process data gathered and calls UI methods to display stats. A StatsManager is launched on a new thread for every stats configuration and at every interval it will collect measures from the WebsiteManagers (which have references to WebsiteStats that can be accessed by stats configuration).

As said before, new incoming measures are first processed by the WebsiteManager, then they will be forwarded to WebsiteStatsHandlers for stats processing but also to the AvailabilityAlert for updating the alert data and checking if a new alert message should be displayed. An alert handler works similarly as a WebsiteStatsHandler because it contains one or more aggregators (depending on the nature of stats) and will purge outdated measures and push a new measure at every update. However, AvailabilityAlert must take care to send the user a message only when the availability is below a certain threshold.

Regarding the UI, the gocui library has been used to create a simple console user interface with scrollable views that can serve the purpose of the application. The UI is adapted depending on the number of stats configuration that were present in the config file.

## Improvements
First of all, the app could be improved in several aspects in terms of documentation and testing. Unfortunately, for lack of time, I couldn't do more than this but it is definitely important to add more comments to explain each part of the project and increase the coverage by having more meaningful tests.

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
