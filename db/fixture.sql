INSERT INTO therapists
  (email, type, name, street_address, city, postcode, country,
   phone, website, languages, short_profile, full_profile,
   status, created_at)
VALUES
  ('test1@example.com', 'ergo', 'User Test1', '111 Test St.', 'Testville', '1111', 'AT',
   '+43123456789', 'https://test1.test', '{"de", "en"}',
   'This is test therapist #1',
   'Test therapist #1 is the first test therapist and is an OT.',
   'approved',
   '2020-04-13 17:27:00'),
  ('test2@example.com', 'physio', 'User Test2', '112 Test St.', 'Testville', '1111', 'AT',
   '+43123456788', 'https://test2.test', '{"de"}',
   'This is test therapist #2',
   'Test therapist #1 is the second test therapist and is a physio.',
   'new',
   '2020-04-13 17:28:00'),
  ('test3@example.com', 'logo', 'User Test3', '113 Test St.', 'Testville', '1111', 'AT',
   '+43123456787', 'https://test3.test', '{"de", "it"}',
   'This is test therapist #3',
   'Test therapist #1 is the third test therapist and is a speech therapist.',
   'edits_pending',
   '2020-04-13 17:29:00'),
  ('test4@example.com', 'ergo', 'User Test4', '114 Test St.', 'Testville', '1111', 'AT',
   '+43123456786', 'https://test4.test', '{"de"}',
   'This is test therapist #4',
   'Test therapist #1 is the fourth test therapist and is an OT.',
   'approved',
   '2020-04-13 17:30:00'),
  ('test5@example.com', 'physio', 'User Test5', '115 Test St.', 'Testville', '1111', 'AT',
   '+43123456785', 'https://test5.test', '{"de", "fr"}',
   'This is test therapist #5',
   'Test therapist #1 is the fifth test therapist and is a physio.',
   'approved',
   '2020-04-13 17:31:00'),
  ('test6@example.com', 'logo', 'User Test6', '116 Test St.', 'Testville', '1111', 'AT',
   '+43123456784', 'https://test6.test', '{"de", "en", "it"}',
   'This is test therapist #6',
   'Test therapist #1 is the sixth test therapist and is a speech therapist.',
   'approved',
   '2020-04-13 17:32:00');
