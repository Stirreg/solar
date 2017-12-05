#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys 
from pymodbus.client.sync import ModbusSerialClient as ModbusClient 
from pymongo import MongoClient 
from datetime import datetime 

modbus = ModbusClient(method='rtu', port='/dev/ttyUSB0', baudrate=9600, stopbits=1, parity='N', bytesize=8, timeout=1) 
modbus.connect() 
read_registers = modbus.read_input_registers(0,33) 
registers = read_registers.registers
#Status Inverter (Stat) 0=waiting 1=normaal 2=fault
status = registers[0] 
if status != 1:
	sys.exit() 

pv_volts = float(registers[3])/10 #1V 
pv_amps = float(registers[4])/10 #1A 
pv_watts = float(registers[6])/10 #1W 
ac_watts = float(registers[12])/10 #1W 
ac_herz = float(registers[13])/100 #1Hz 
ac_volts = float(registers[14])/10 #1V 
ac_amps = float(registers[15])/10 #1A 
total_today = float(registers[27])/10 #1kWh 
total = float(registers[29])/10 #1kWh
#Total runtime (Htot-low byte) 0.5S
runtime_low = float(registers[30])
#Total runtime (Htot-high-byte) 0.5S
runtime_high = float(registers[31]) 
temperature = float(registers[32])/10 #1C 

modbus.close() 

mongo = MongoClient() 
db = mongo.pv 

db.growatt.insert(
	{
        	'created_at': datetime.utcnow(),
		'status': status,
		'pv_volts': pv_volts,
		'pv_amps': pv_amps,
		'pv_watts': pv_watts,
		'ac_watts': ac_watts,
		'ac_herz': ac_herz,
		'ac_volts': ac_volts,
		'ac_amps': ac_amps,
		'total_today': total_today,
		'total': total,
		'runtime_low': runtime_low,
		'runtime_high': runtime_high,
		'temperature': temperature
	}
)
		
		
		
