#!/bin/bash

# Supply Chain - Bahan Baku
cat > test-cases/supply-chain-bahan-baku/test-cases.json << 'EOF'
[
  {"id":"rm-001","module":"supply-chain-bahan-baku","scenario":"View raw material inventory","steps":["Login","Navigate to Bahan Baku","Wait for inventory","Verify list"],"expectedResults":["Inventory displayed","Stock levels shown","Material details visible","Units shown"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"rm-002","module":"supply-chain-bahan-baku","scenario":"Add new raw material","steps":["Navigate to Bahan Baku","Click add","Enter material details","Set initial stock","Save"],"expectedResults":["Form displayed","Material created","Success message","Material in list"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"rm-003","module":"supply-chain-bahan-baku","scenario":"Update stock levels","steps":["Navigate to Bahan Baku","Select material","Update quantity","Save changes"],"expectedResults":["Stock updated","History recorded","Alert if low stock","Success message"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"rm-004","module":"supply-chain-bahan-baku","scenario":"Low stock alert","steps":["Navigate to Bahan Baku","View materials","Check low stock items","Verify alerts"],"expectedResults":["Low stock highlighted","Alert shown","Reorder suggestion","Notification sent"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]}
]
EOF

# Logistik - Data Sekolah
cat > test-cases/logistik-data-sekolah/test-cases.json << 'EOF'
[
  {"id":"sch-001","module":"logistik-data-sekolah","scenario":"View school data","steps":["Login","Navigate to Data Sekolah","Wait for schools","Verify list"],"expectedResults":["School page displayed","List loaded","School details visible","Location shown"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"sch-002","module":"logistik-data-sekolah","scenario":"Create new school record","steps":["Navigate to schools","Click add","Enter school details","Set location","Save"],"expectedResults":["Form displayed","School created","Success message","School in list"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"sch-003","module":"logistik-data-sekolah","scenario":"Edit school information","steps":["Navigate to schools","Select school","Click edit","Modify details","Save"],"expectedResults":["Edit form shown","Changes saved","Updated data displayed","History logged"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"sch-004","module":"logistik-data-sekolah","scenario":"Manage school location","steps":["Navigate to schools","Select school","Update location","Set coordinates","Save"],"expectedResults":["Location form shown","Coordinates saved","Map updated","Location verified"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["medium"]}
]
EOF

# Logistik - Tugas Pengiriman
cat > test-cases/logistik-tugas-pengiriman/test-cases.json << 'EOF'
[
  {"id":"del-001","module":"logistik-tugas-pengiriman","scenario":"View delivery tasks","steps":["Login","Navigate to Tugas Pengiriman","Wait for tasks","Verify list"],"expectedResults":["Task page displayed","List loaded","Task details visible","Status shown"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]},
  {"id":"del-002","module":"logistik-tugas-pengiriman","scenario":"Create delivery task","steps":["Navigate to tasks","Click create","Select school","Set delivery date","Assign driver","Save"],"expectedResults":["Form displayed","Task created","Success message","Task in list"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"del-003","module":"logistik-tugas-pengiriman","scenario":"Update task status","steps":["Navigate to tasks","Select task","Change status","Add notes","Save"],"expectedResults":["Status updated","Timestamp recorded","Notification sent","History logged"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["critical"]},
  {"id":"del-004","module":"logistik-tugas-pengiriman","scenario":"Assign pickup task","steps":["Navigate to tasks","Create pickup","Select location","Assign driver","Save"],"expectedResults":["Pickup created","Driver notified","Task scheduled","Success message"],"actualResults":[],"status":"not_run","lastExecuted":null,"executionTime":null,"tags":["high"]}
]
EOF

echo "Test cases created successfully"
