#!/usr/bin/python3

import os
import re
import csv
import pathlib

output_file_name="findings.txt"

class Helper:
	def __init__(self):
		this_file=os.path.abspath(__file__)
		self.this_file=this_file[:1].lower()+this_file[1:]		
		
	def check_dir_path(self, path):
		if(not os.path.exists(path) or not os.path.isdir(path)):
			print('Invalid directory path specified')
			exit()
			
	def check_file_path(self, path):
		if(not os.path.exists(path)):
			print('Invalid file path specified')
			exit()
	
	def header(self):
		print("\r\nThank you, working with data now...\r\n")
	
	def footer(self,counter,output_path):
		print("\r\n"+str(counter)+" matches have been found."+"\r\nThey are stored in "+output_path+".")
		print("\r\nHave a nice day !\r\n")

if __name__=='__main__':			
	helper=Helper();
	data_path=input("Please specify absolute path to the data directory:\r\n")
	helper.check_dir_path(data_path);
	output_path=input("Please specify absolute path to the output directory:\r\n")
	helper.check_dir_path(output_path);
	#files=os.listdir(data_path)
	path=pathlib.Path(data_path)
	files_and_dirs=list(path.rglob('*'))
	files=[str(file) for file in files_and_dirs if file.is_file()];

	pattern_path=input("Please specify absolute path to the file with patterns:\r\n")
	helper.check_file_path(pattern_path);#print(pattern_path);
	
	patterns_file=open(pattern_path,'r')
	patterns=patterns_file.readlines();
	patterns_file.close()

	# Filter out current file and search patterns if in working dir
	files=list(filter(lambda file: file != pattern_path and file != helper.this_file, files));#print(files)
	files.sort();
	
	counter=0

	patterns=[pat.rstrip('\r\n') for pat in patterns];
	pats_to_print=patterns
	patterns=['.*'+pat+'.*' for pat in patterns]
	patterns=[re.compile(pat) for pat in patterns]

	output_file=open(output_path+'\\'+output_file_name,'w',newline='')
	writer=csv.writer(output_file)
	
	helper.header()
	
	for file_name in files:
		#cur_file_name=data_path+'\\'+file_name;#print(cur_file_name);	
		with open(file_name,'r',1,errors='ignore') as file:
			print('Searching '+file_name+'...')
			found=0
			for line in file:
				cur_pat=0
				for pattern in patterns:
					if(pattern.match(line)): 
						last_slash=file_name.rfind('\\')
						result=[file_name[last_slash+1:],file_name[:last_slash],line.rstrip('\r\n'),pattern.pattern[2:-2]];#print(result)
						print('	'+pats_to_print[cur_pat])
						#file=open(output_path+'\\'+'log'+'-'+str(counter)+'.txt','w')
						#writer=csv.writer(file)
						writer.writerow(result)					
						counter+=1
						found=1
					cur_pat+=1
			if not found: print('	No matches found.')
	output_file.close()
	helper.footer(counter,output_path)