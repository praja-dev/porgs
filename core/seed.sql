WITH cts AS (SELECT strftime('%s', 'now') AS now)
INSERT
INTO org_type (id, created, updated, name, description, cxp)
VALUES (1, (SELECT now FROM cts), (SELECT now FROM cts), 'Org', 'Generic Organization', ''),
       (2, (SELECT now FROM cts), (SELECT now FROM cts), 'Country', 'Generic Country', ''),
       (3, (SELECT now FROM cts), (SELECT now FROM cts), 'Region', 'Generic Region', ''),
       (4, (SELECT now FROM cts), (SELECT now FROM cts), 'State', 'Generic State', ''),
       (5, (SELECT now FROM cts), (SELECT now FROM cts), 'Province', 'Generic Province', ''),
       (6, (SELECT now FROM cts), (SELECT now FROM cts), 'County', 'Generic County', ''),
       (7, (SELECT now FROM cts), (SELECT now FROM cts), 'District', 'Generic District', ''),
       (8, (SELECT now FROM cts), (SELECT now FROM cts), 'City', 'Generic City', ''),
       (9, (SELECT now FROM cts), (SELECT now FROM cts), 'Town', 'Generic Town', ''),
       (10, (SELECT now FROM cts), (SELECT now FROM cts), 'Village', 'Generic Village', ''),
       (1000, (SELECT now FROM cts), (SELECT now FROM cts), 'Sri Lanka', 'The country of Sri Lanka', ''),
       (1001, (SELECT now FROM cts), (SELECT now FROM cts), 'Province', 'A province in Sri Lanka', ''),
       (1002, (SELECT now FROM cts), (SELECT now FROM cts), 'District', 'A district in Sri Lanka', ''),
       (1003, (SELECT now FROM cts), (SELECT now FROM cts), 'DS Division', 'A DS Division in Sri Lanka', ''),
       (1004, (SELECT now FROM cts), (SELECT now FROM cts), 'GN Division', 'A GN Division in Sri Lanka', ''),
       (1005, (SELECT now FROM cts), (SELECT now FROM cts), 'Village', 'A village in Sri Lanka', ''),
       (1010, (SELECT now FROM cts), (SELECT now FROM cts), 'Electoral District', 'An electoral district in Sri Lanka',
        ''),
       (1011, (SELECT now FROM cts), (SELECT now FROM cts), 'Polling Division', 'A polling division in Sri Lanka', ''),
       (1012, (SELECT now FROM cts), (SELECT now FROM cts), 'Polling District', 'A polling district in Sri Lanka', '');



