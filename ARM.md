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

<li>  </li>

<li>  </li>

<li>  </li>

<li>  </li>

<li>  </li>

<li>  </li>

</ul>
