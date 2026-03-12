#!/bin/bash

# SDM - Data Karyawan
cat > test-cases/sdm-data-karyawan/test-cases.json << 'EOF'
[
  {"id":"emp-001","module":"sdm-data-karyawan","scenario":"View employee list","steps":["Login","Navigate to Data Karyawan","Wait for employees","Verify list"],"expectedResults":["Employee page displayed","List loaded","Employee details visible","Roles shown"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"emp-002","module":"sdm-data-karyawan","scenario":"Create new employee","steps":["Navigate to employees","Click add","Enter employee details","Set role","Save"],"expectedResults":["Form displayed","Employee created","Success message","Employee in list"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"emp-003","module":"sdm-data-karyawan","scenario":"Edit employee information","steps":["Navigate to employees","Select employee","Click edit","Modify details","Save"],"expectedResults":["Edit form shown","Changes saved","Updated data displayed","History logged"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"emp-004","module":"sdm-data-karyawan","scenario":"Manage employee roles","steps":["Navigate to employees","Select employee","Update role","Set permissions","Save"],"expectedResults":["Role form shown","Role updated","Permissions applied","Access changed"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]}
]
EOF

# SDM - Laporan Absensi
cat > test-cases/sdm-laporan-absensi/test-cases.json << 'EOF'
[
  {"id":"att-001","module":"sdm-laporan-absensi","scenario":"View attendance report","steps":["Login","Navigate to Laporan Absensi","Select date range","View report"],"expectedResults":["Report page displayed","Attendance data shown","Statistics visible","Export available"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"att-002","module":"sdm-laporan-absensi","scenario":"Filter by employee","steps":["Navigate to report","Select employee filter","Choose employee","Apply filter"],"expectedResults":["Filter applied","Employee data shown","Attendance records filtered","Summary updated"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"att-003","module":"sdm-laporan-absensi","scenario":"Filter by date range","steps":["Navigate to report","Select date filter","Choose range","Apply filter"],"expectedResults":["Date filter applied","Data filtered","Statistics updated","Report refreshed"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"att-004","module":"sdm-laporan-absensi","scenario":"Export attendance report","steps":["Navigate to report","Configure filters","Click export","Download file"],"expectedResults":["Export dialog shown","File generated","Download started","File format correct"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["medium"]}
]
EOF

# SDM - Konfigurasi Absensi
cat > test-cases/sdm-konfigurasi-absensi/test-cases.json << 'EOF'
[
  {"id":"atc-001","module":"sdm-konfigurasi-absensi","scenario":"View attendance configuration","steps":["Login","Navigate to Konfigurasi Absensi","View settings","Verify config"],"expectedResults":["Config page displayed","Settings shown","Work hours visible","Rules displayed"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"atc-002","module":"sdm-konfigurasi-absensi","scenario":"Set work hours","steps":["Navigate to config","Click edit hours","Set start/end time","Save"],"expectedResults":["Hours form shown","Times saved","Config updated","Success message"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"atc-003","module":"sdm-konfigurasi-absensi","scenario":"Configure attendance rules","steps":["Navigate to config","Edit rules","Set late penalty","Save rules"],"expectedResults":["Rules form shown","Rules saved","Penalties configured","Applied to system"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"atc-004","module":"sdm-konfigurasi-absensi","scenario":"Setup IP-based check-in","steps":["Navigate to config","Enable IP check","Add allowed IPs","Save config"],"expectedResults":["IP config shown","IPs saved","Validation enabled","Check-in restricted"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]}
]
EOF

# Keuangan - Aset Dapur
cat > test-cases/keuangan-aset-dapur/test-cases.json << 'EOF'
[
  {"id":"ast-001","module":"keuangan-aset-dapur","scenario":"View kitchen assets","steps":["Login","Navigate to Aset Dapur","Wait for assets","Verify list"],"expectedResults":["Assets page displayed","List loaded","Asset details visible","Values shown"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"ast-002","module":"keuangan-aset-dapur","scenario":"Add new asset","steps":["Navigate to assets","Click add","Enter asset details","Set value","Save"],"expectedResults":["Form displayed","Asset created","Success message","Asset in list"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"ast-003","module":"keuangan-aset-dapur","scenario":"Track asset depreciation","steps":["Navigate to assets","Select asset","View depreciation","Check value"],"expectedResults":["Depreciation shown","Value calculated","History displayed","Report available"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"ast-004","module":"keuangan-aset-dapur","scenario":"Record maintenance","steps":["Navigate to assets","Select asset","Add maintenance record","Save"],"expectedResults":["Maintenance form shown","Record saved","History updated","Cost recorded"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["medium"]}
]
EOF

# Keuangan - Arus Kas
cat > test-cases/keuangan-arus-kas/test-cases.json << 'EOF'
[
  {"id":"cf-001","module":"keuangan-arus-kas","scenario":"View cash flow records","steps":["Login","Navigate to Arus Kas","Wait for records","Verify list"],"expectedResults":["Cash flow page displayed","Records loaded","Inflows/outflows shown","Balance visible"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"cf-002","module":"keuangan-arus-kas","scenario":"Record cash inflow","steps":["Navigate to cash flow","Click add inflow","Enter amount","Set category","Save"],"expectedResults":["Inflow form shown","Record created","Balance updated","Success message"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"cf-003","module":"keuangan-arus-kas","scenario":"Record cash outflow","steps":["Navigate to cash flow","Click add outflow","Enter amount","Set category","Save"],"expectedResults":["Outflow form shown","Record created","Balance updated","Success message"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"cf-004","module":"keuangan-arus-kas","scenario":"View cash flow report","steps":["Navigate to cash flow","Select date range","Generate report","View summary"],"expectedResults":["Report generated","Summary shown","Charts displayed","Export available"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]}
]
EOF

# Keuangan - Laporan
cat > test-cases/keuangan-laporan/test-cases.json << 'EOF'
[
  {"id":"fin-001","module":"keuangan-laporan","scenario":"View financial reports","steps":["Login","Navigate to Laporan Keuangan","Select report type","View report"],"expectedResults":["Report page displayed","Report loaded","Data shown","Charts visible"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"fin-002","module":"keuangan-laporan","scenario":"Generate profit/loss statement","steps":["Navigate to reports","Select P&L","Set date range","Generate"],"expectedResults":["P&L generated","Revenue shown","Expenses shown","Profit calculated"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"fin-003","module":"keuangan-laporan","scenario":"Generate balance sheet","steps":["Navigate to reports","Select balance sheet","Set date","Generate"],"expectedResults":["Balance sheet generated","Assets shown","Liabilities shown","Equity calculated"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"fin-004","module":"keuangan-laporan","scenario":"Export financial report","steps":["Navigate to reports","Generate report","Click export","Download file"],"expectedResults":["Export dialog shown","File generated","Download started","Format correct"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["medium"]}
]
EOF

# Sistem - Audit Trail
cat > test-cases/sistem-audit-trail/test-cases.json << 'EOF'
[
  {"id":"aud-001","module":"sistem-audit-trail","scenario":"View audit logs","steps":["Login","Navigate to Audit Trail","Wait for logs","Verify list"],"expectedResults":["Audit page displayed","Logs loaded","Actions shown","Users visible"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"aud-002","module":"sistem-audit-trail","scenario":"Filter by user","steps":["Navigate to audit","Select user filter","Choose user","Apply filter"],"expectedResults":["Filter applied","User logs shown","Actions filtered","Timeline displayed"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"aud-003","module":"sistem-audit-trail","scenario":"Filter by action type","steps":["Navigate to audit","Select action filter","Choose type","Apply filter"],"expectedResults":["Filter applied","Actions filtered","Logs shown","Count updated"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["medium"]},
  {"id":"aud-004","module":"sistem-audit-trail","scenario":"Export audit logs","steps":["Navigate to audit","Configure filters","Click export","Download"],"expectedResults":["Export dialog shown","File generated","Download started","Format correct"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["medium"]}
]
EOF

# Sistem - Konfigurasi
cat > test-cases/sistem-konfigurasi/test-cases.json << 'EOF'
[
  {"id":"sys-001","module":"sistem-konfigurasi","scenario":"View system configuration","steps":["Login","Navigate to Konfigurasi","View settings","Verify config"],"expectedResults":["Config page displayed","Settings shown","Parameters visible","Values displayed"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"sys-002","module":"sistem-konfigurasi","scenario":"Update system settings","steps":["Navigate to config","Edit settings","Modify values","Save"],"expectedResults":["Settings form shown","Values saved","Config updated","Success message"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"sys-003","module":"sistem-konfigurasi","scenario":"Manage user permissions","steps":["Navigate to config","Select permissions","Modify roles","Save"],"expectedResults":["Permission form shown","Roles updated","Access changed","Applied immediately"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"sys-004","module":"sistem-konfigurasi","scenario":"Backup and restore settings","steps":["Navigate to config","Click backup","Download config","Verify backup"],"expectedResults":["Backup created","File downloaded","Settings saved","Restore available"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]}
]
EOF

echo "All remaining test cases created successfully"
