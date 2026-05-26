# shift-scheduler
The capstone assignment for OIT Full Stack Training. This repository contains the code for a shift scheduler similar to the one used by OIT.

To run, ensure that you have golang installed on your machine. Navigate to the home directory, then run the following command:

```console
go run .
```

Once the server is started, go to http://localhost:4000/schedule to view the basic schedule view. Click and drag to create the schedule, then submit it for admin approval. Click the button in the upper right to switch from the student to the admin view.

Once in the admin view, an approval tab will appear. User/Schedule 0 is the student schedule, and User/Schedule 1 is the admin schedule. Approve or reject the schedule from there. If rejecting, include an optional rejection message. Then switch back to the student view to see that the schedule is approved or rejected. If rejected, edit it further and resubmit it. You can continue to follow that process until the schedule is approved.
