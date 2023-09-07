example rows in pto approval worksheet:

```
[Request Details]
[Date Day of the Week Type Requested Unit of Time]
[09-15-23 Freitag Annual Leave 1 Days]
[09-29-23 Freitag Annual Leave 1 Days]
[10-13-23 Freitag Annual Leave 1 Days]
[10-27-23 Freitag Annual Leave 1 Days]
[11-10-23 Freitag Annual Leave 1 Days]
[11-13-23 Montag Annual Leave 1 Days]
[11-14-23 Dienstag Annual Leave 1 Days]
[11-24-23 Freitag Annual Leave 1 Days]
[12-08-23 Freitag Annual Leave 1 Days]
[12-22-23 Freitag Annual Leave 1 Days]
```

usage:

```
go run ./cmd/workday-excel-to-ical \
  --organiser jgwosdz@redhat.com \
  --domain pto.sd.redhat.com \
  --summary "jgwosdz - PTO" \
  ~/Downloads/Absence_Request.xlsx > output.ical
```
