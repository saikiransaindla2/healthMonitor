# Health Monitor

Health Monitor is a small service created to monitor the health of external URLs




## Working

System takes in a list of http/https urls and a crawlTime(seconds) and waitTime(in seconds) and threshold(count)
where:

crawlTime : System will wait for this much time before giving up on the url
waitTime : System will wait for this much time before retrying again.
threshold : Count of retries possible for that url

The system shall iterate over all the urls in the system and try to do a HTTP GET on the URL(wait for the crawl_timeout) seconds before giving up on the URL.

The health check status is inserted in the table testing_data and can be displayed through the route "/HealthMonitor/fetch/{runId}" for a particular runId
