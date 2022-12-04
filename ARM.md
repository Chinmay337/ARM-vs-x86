<link href="style.css" rel="stylesheet" />

# ARM v86

<ul>

<li> AWS Scheduler </li>

</br>

```
First create periods

This creates a period that will start/stop our EC2's 3 times a day

scheduler-cli create-period --stack ARM-Scheduler --name armmorn --begintime 09:00 --endtime 09:30 --region us-east-1

scheduler-cli create-period --stack ARM-Scheduler --name armeve --begintime 15:00 --endtime 15:30 --region us-east-1

scheduler-cli create-period --stack ARM-Scheduler --name armnight --begintime 22:00 --endtime 22:30 --region us-east-1

Then add periods to a schedule

scheduler-cli create-schedule --name ARMSchedule --periods armmorn,armeve,armnight --stack ARM-Scheduler --region us-east-1

Deleting schedule
scheduler-cli delete-schedule --name ARMSchedule --stack ARM-Scheduler --region us-east-1
```

</br>

<li> Querying Usage</li>

</br>

```
ARM
parse @message /Duration:\s*(?<@duration_ms>\d+\.\d+)\s*ms\s*Billed\s*Duration:\s*(?<@billed_duration_ms>\d+)\s*ms\s*Memory\s*Size:\s*(?<@memory_size_mb>\d+)\s*MB/
| filter @message like /REPORT RequestId/
| stats sum(@billed_duration_ms * @memory_size_mb * 1.33334e-11 + 2.0e-7) as @cost_dollars_total

x86
parse @message /Duration:\s*(?<@duration_ms>\d+\.\d+)\s*ms\s*Billed\s*Duration:\s*(?<@billed_duration_ms>\d+)\s*ms\s*Memory\s*Size:\s*(?<@memory_size_mb>\d+)\s*MB/
| filter @message like /REPORT RequestId/
| stats sum(@billed_duration_ms * @memory_size_mb * 1.66667e-11 + 2.0e-7) as @cost_dollars_total
```

</br>

<li> Results </li>

</br>

```
ARM Go Float
@cost_dollars_total 0.00877

Intel Go Float
@cost_dollars_total	0.007329

Intel Go Int
@cost_dollars_total	0.03283

ARM Go Int
@cost_dollars_total	0.01104

ARM JS Int
@cost_dollars_total	0.01434 1m 0.03659

Intel JS Int
@cost_dollars_total	0.01437 1m 0.03658


ARM JS Float
@cost_dollars_total	0.01523

Intel JS Float
@cost_dollars_total	0.01269


```

</br>

<li>  </li>

<li>  </li>

<li>  </li>

<li>  </li>

</ul>
