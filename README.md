# PTSD

Pager Trauma Statistics Daemon calls out to PagerDuty every so often to grab
statistics on:

 * total number of incidents 
 * notifications per user
 * assignments per user
 * acknowledgements per user

It then sends these as counter metrics to Statsd. It was created to be modular
with respect to collectors and outputters so support for alternatives like
VictorOps are coming soon (feel free to PR me).

It's a Go program, configured via environment variables, so you can run it
locally, in a Container, on Heroku, or pretty much anywhere else.

Current Environment Variables:

 * PTSD_INTERVAL: (defaults to 60) specifies (in minutes) how often PTSD polls
	each collector for metrics.
 * PDTOKEN: enables PagerDuty collector and specifies the API auth-token to use.
 * STATSD_SERVER & STATSD_PORT: Enables Statsd outputter, and specifies the server IP and port (these must both be set to enable statsd)
 * PTSD_TXT: enables metrics output to STDOUT (set to '1' or anything else to enable)
 * DEBUG: enables debugging output to STDOUT (set to '1' or anything else to
	enable debugging)

## Re: Configuration
A couple things you should know about configuring PTSD: 

 * The collectors/outputters are not mutually exclusive. So if you enable
	multiple collectors they will all be called every *interval* minutes
 * The collector/outputter Environment vars are parsed every interval, so you
	don't need to bring down PTSD to enable/disable collectors/outputters or
	change your PD credentials etc..

## Extending

PTSD is currently neither godoc'd nor gofmt'd. I don't really feel bad about
either of those things since I mostly wrote this on the plane-ride home from
monitorama.

But if you speak Go and you'd like to extend PTSD to support alternative
collectors or outputters just drop a file in the root dir that implements the
`Outputter` or `Collector` interface (see [interfaces.go](/interfaces.go)).
Then append your new thingy to the global `OUTPUTTERS` or `COLLECTORS` slice
[like so](https://github.com/djosephsen/ptsd/blob/master/pagerduty.go#L24-L26)
