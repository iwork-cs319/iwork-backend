INSERT INTO offerings(id, workspace_id, user_id, start_time, end_time, cancelled)
VALUES ('94a019ec-22ef-41b1-979c-3c52567dfcba', 'c534e44e-a80c-4f26-97b7-1c8ae7f04856', 'd2e36d1b-d944-4f79-b263-2e53ff77dcc3', '2019-10-14 00:00:00.000000', '2019-10-15 23:59:59.000000', false)
RETURNING id;

INSERT INTO offerings(id, workspace_id, user_id, start_time, end_time, cancelled)
VALUES ('8cdb5375-a2bc-4a82-9862-b37bdd919265', '151d54fd-343b-48fb-b90a-8a68aafed512', '96e99446-22a4-4d54-9c81-dc54a894781d', '2019-10-11 00:00:00.000000', '2019-10-15 23:59:59.000000', false)
RETURNING id;