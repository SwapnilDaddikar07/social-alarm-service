# social-alarm-service

Create Alarm Service.

The app allows you to schedule alarms on your device.
These alarms can be seen by your contacts.
Your contacts can decide to send you an audio/video which will be played on your device when the alarm goes off.

APIs available as of now.
More APIs will be added as the development progresses.

1. Create alarm
   Allows users to create either public or private alarms. Users can add description about the alarm.

2. Update alarm status
   Allows users to turn an alarm ON/OFF.

3. Eligible alarms
   When called with a user-id , this alarm returns public , non-expired repeating and non-repeating alarms.

4. All alarms
   Used to fetch all alarms (public/private or expired/non-expired). This API will be called by the owner of the alarms.

5. Upload media
   This API allows your contacts to record media (audio/media) for an alarm you've set.

6. Fetch media
   This API is called when an alarm goes off. This API fetches all available media links for an alarm id. Frontend will play the media associated with these links.