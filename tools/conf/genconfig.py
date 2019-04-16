#!/usr/bin/python
# -*- coding:utf-8 -*-
import os
parseList = {"bool":"antnet.ParseBool","int":"antnet.Atoi","int32":"antnet.ParseInt32","uint32":"antnet.ParseUint32","int64":"antnet.Atoi64","uint64":"antnet.ParseUint64","float32":"antnet.Atof","float64":"antnet.Atof64"}
baseTypeList = {"boolean","long","double","bool","int","int32","uint32","int64","uint64","float","float32","float64","string"}
type_map = {"long":"int64","float":"float32","double":"float64","boolean":"bool"}
typeList = ['int','float','double','string']
rowDef = {'type':0,"commit":1,'key':2,'ParseType':3}

def ListDirFile(rootdir):
	files = []
	list = os.listdir(rootdir) #列出文件夹下所有的目录与文件
	for i in range(0,len(list)):
		path = os.path.join(rootdir,list[i])
		if os.path.isfile(path):
			files.append(path)
	return files


#解析字段1 服务器 2 客户端 3客户端和服务器公用
def canParse(type):
	if type == "1" or type == "3":
		return True
	return False

#自定义类型的解析
def parseFixedListFun(type):
	parseType = ""
	if type == "string":
		parseType = "val" 
	else:
		parseType = "{0}(val)".format(parseList[type])

	funcName = "parse{0}ListFun".format(type[0].upper()+type[1:])
	str="""
var %s = func(fieldv reflect.Value, data, path string) error{
	if data == "0" || data == ""{
		return nil
	}
	sl1 := antnet.StrSplit(data,";")
	var dataList = make([]%s,0,len(sl1))
	for _, val:= range sl1{
		dataList = append(dataList,%s)
	}
	fieldv.Set(reflect.ValueOf(dataList))
	return nil
}"""%(funcName,type,parseType)
	return funcName,str

#通用配置excel
def genConfig(fname,exportStructInfo,fixedDefInfo):
	configStructStr=""
	configStructReadStr=""
	configStructGetStr="" 
	configFieldStr=""
	readConfName=""

	lines = []
	(filepath,tempfilename) = os.path.split(fname)
	(shotname,extension) = os.path.splitext(tempfilename)
	structName = shotname[0].upper() + shotname[1:]
	with open(fname,"r") as f:
		old=f.readlines()
	for l in old:
		s = l.split('\t')
		ss = []
		for sss in s:
			ss.append(sss.strip())
		lines.append(ss)
	if len(old) < 5:
		return "","","","",""
	keys = []
	Id= "Id"
	idType = 'int'
	structFieldString = ""
	hasSelfDefType = False

	#print("structName:{0},line:{1}".format(structName, lines[0]))
	for i, key in enumerate(lines[rowDef['key']]):
		if key == 'id':
			idType = lines[rowDef['type']][i]
		if structName == "Global" and key == "key":
			idType = lines[rowDef['type']][i]
			Id = "Key"
		parsetype = lines[rowDef['ParseType']][i]
		if canParse(parsetype) == True:
			keyType =  lines[rowDef['type']][i]
			if keyType in baseTypeList:
				if keyType in type_map:
					keyType = type_map[keyType]
			elif keyType.endswith("[]"):
				hasSelfDefType =True
				tmpType = keyType[:-2]
				if tmpType in exportStructInfo:
					keyType = "[]"+tmpType
					info = exportStructInfo[tmpType]
					info[3] = True
					exportStructInfo[tmpType] = info
				elif tmpType in baseTypeList:
					if tmpType in type_map:
						tmpType = type_map[tmpType]
					if tmpType not in fixedDefInfo:
						funcName,funStr = parseFixedListFun(tmpType)
						fixedDefInfo[tmpType] = [funcName,funStr]
					keyType = "[]"+tmpType
				else:
					print("error: parse {0} can not find type {1}".format(structName,tmpType))
					os.system("pause")
					return "","","","",""
			elif keyType in exportStructInfo:
				hasSelfDefType = True
				info = exportStructInfo[keyType]
				info[0] = True
				exportStructInfo[keyType] = info
				#print("excel:{0},fieldType:{1}".format(structName,keyType))
			else:
				print("error: can not parse excel:{0} row:1, line:{1} type:{2}".format(structName,i+1,keyType))
				os.system("pause")
				return "","","","",""
			newKey = key[0].upper() + key[1:]
			structFieldString+='\n\t%s\t%s `json:"%s"`' %(newKey, keyType, key)
			keys.append([newKey,keyType])
	if len(keys) == False:
		return "","","","",""

	selfDefStr =""
	parseObjFunStr = ""
	if hasSelfDefType:
			selfDefStr+="\n\tparseObjFun := make(map[string]func(fieldv reflect.Value, data, path string)error)"
			parseObjFunStr = "ParseObjFun : parseObjFun,"
			for k in keys:
				if k[1].startswith("[]"):
					tmpType = k[1][2:]
					if tmpType in fixedDefInfo:
						selfDefStr+= "\n\tparseObjFun[\"{0}\"] ={1}".format(k[0],fixedDefInfo[tmpType][0])
					elif tmpType in exportStructInfo:
						selfDefStr+= "\n\tparseObjFun[\"{0}\"] ={1}".format(k[0],exportStructInfo[tmpType][4])
				elif k[1] in exportStructInfo:
					selfDefStr += "\n\tparseObjFun[\"{0}\"] ={1}".format(k[0],exportStructInfo[k[1]][1])

	configStructStr='''
type Config%s struct {%s
}'''%(structName,structFieldString)


	ReadConfigName = "ReadConfig{0}".format(structName)
	configStructReadStr='''
func %s() map[%s]*Config%s {%s
	m := map[%s]*Config%s{}
	err, datas := antnet.ReadConfigFromCSV("conf/doc/%s", 3, 5, &antnet.GenConfigObj{
		GenObjFun: func() interface{} {
			return &Config%s{}
		},%s
	})
	if err == nil{
		for _, data := range datas{
			config := data.(*Config%s)
			m[config.%s] = config
		}
		antnet.LogInfo("read conf success path:conf/doc/%s")
	} else{
		antnet.LogError("read conf failed path:conf/doc/%s err:%%v", err)
	}
	return m
}'''%(ReadConfigName,idType, structName,selfDefStr,idType,structName, tempfilename, structName, parseObjFunStr,structName, Id,tempfilename,tempfilename)


	configStructGetStr='''
func (this *configDoc) GetConfig%s(id %s) *Config%s {
	if value, founded := this.%s[id]; founded{
		return value
	}
	return nil
}''' %(structName,idType,structName,structName)


	configFieldStr = "\n\t{0} map[{1}]*Config{0} `json:\"{0}\"`".format(structName,idType)
	readConfName  = "\n\t\t{0}(),".format(ReadConfigName)

	return configStructStr, configStructReadStr, configStructGetStr,configFieldStr,readConfName



#自定义数据
def parseSelfDefFun(structName,fieldNames,fieldTypes):
	filedStr = ""
	fieldNum = len(fieldNames)
	for index, field in enumerate(fieldNames):
		type = fieldTypes[index]
		if type == "string":
			filedStr += "\n\ttmpValue.%s = sl1[%s]"%(field,index) 
		else:
			filedStr += "\n\ttmpValue.%s = %s(sl1[%s])"%(field,parseList[type],index) 
	funcName = "parse{0}Fun".format(structName)
	str = """
var	%s = func(fieldv reflect.Value, data, path string) error{
	if data == "0" || data == ""{
		return nil
	}
	var sl1 []string
	antnet.SplitString1(data,&sl1)
	if len(sl1) < %s{
		return errors.New("field length < %s")
	}
	var tmpValue %s
	fieldv.Set(reflect.ValueOf(tmpValue))
	return nil
}
"""%(funcName,fieldNum,fieldNum,structName+filedStr)
	return funcName, str

#自定义数组数据
def parseSelfDefListFun(structName,fieldNames,fieldTypes):
	
	#自定义的数组类型
	filedStr = ""
	fieldNum = len(fieldNames)
	filedStr+="if len(value) < %s{\n\t\t\treturn errors.New(\"field length < %s\")\n\t\t}"%(fieldNum,fieldNum)
	for index, field in enumerate(fieldNames):
		type = fieldTypes[index]
		if type == "string":
			filedStr += "\n\t\ttmpValue.%s = value[%s]"%(field,index) 
		else:
			filedStr += "\n\t\ttmpValue.%s = %s(value[%s])"%(field,parseList[type],index) 
	funcName = "parse{0}ListFun".format(structName)
	str = """
var %s = func(fieldv reflect.Value, data, path string) error{
	if data == "0"|| data == ""{
		return nil
	}
	var sl2 [][]string
	antnet.SplitString2(data,&sl2)
	var dataList = make([]%s,0,len(sl2))
	var tmpValue %s
	for _, value := range sl2 {
		%s
		dataList = append(dataList,tmpValue)
	}
	fieldv.Set(reflect.ValueOf(dataList))
	return nil
}
"""%(funcName,structName,structName,filedStr)
	return funcName,str



#读取ExportSetting自定义数据类型
def genSelfDefSetting(fname):
	exportStructStr = ""
	exportStructInfo = {}
	with open(fname,"r") as f:
		old=f.readlines()
	if len(old) < 5:
		return exportStructStr, exportStructInfo
	lines = old[4:]
	#f2 = open("%s/frame/auto_config.go" %srcroot,"a+")
	for line in lines:
		s = line.split('\t')
		if len(s)<4:
			print("parse ExportSetting error line info {0}".format(line))
			os.system("pause")
			return exportStructStr,exportStructInfo
		structName = s[0]
		fields = s[2].split(";")
		fieldTypes = s[3].split(";")
		if len(fields) != len(fieldTypes):
			print("parse ExportSetting error struct:{0},fiels length is not equal to types length".format(structName))
			os.system("pause")
			return exportStructStr,exportStructInfo
		newFields = []
		newTypes = []
		fieldStr = ""
		for index, field in enumerate(fields):
			fieldType = fieldTypes[index]
			if fieldType in type_map:
				fieldType = type_map[fieldType]
			field = field[0].upper()+field[1:]
			newFields.append(field)
			newTypes.append(fieldType)
			fieldStr += "\n\t{0} {1}".format(field,fieldType)
		selfDefFunName,selfDefFunStr = parseSelfDefFun(structName,newFields,newTypes)
		selfDefListFunName,selfDefListFunStr = parseSelfDefListFun(structName,newFields,newTypes)
		exportStructInfo[structName] = [False, selfDefFunName,selfDefFunStr,False, selfDefListFunName,selfDefListFunStr]
		exportStructStr+='''
type %s struct {%s
}
''' %(structName,fieldStr)

	return exportStructStr, exportStructInfo


#读取globalSet表格
def genGlobalSetConfig(fname):
	(filepath,tempfilename) = os.path.split(fname)
	keyStr = ""
	structName = "GlobalSet"
	with open(fname,"r") as f:
		old=f.readlines()
	for l in old:
		s = l.split('\t')
		ss = []
		for sss in s:
			ss.append(sss.strip())
		keyStr+= "\n\t{0} {1}  `json:\"%s\"`".format(ss[1][0].upper()+ss[1][1:],ss[2],ss[1])
	str='''
type Config%s struct {%s
}
func ReadConfig%s() *Config%s {
	err, obj := antnet.ReadConfigFromCSVLie("conf/doc/%s", 2, 4, 1, &antnet.GenConfigObj{
		GenObjFun: func() interface{} {
			return &Config%s{}
		},
	})
	if err != nil{
		antnet.LogError("read conf failed path:conf/doc/%s err:%%v", err)
	} else {
		antnet.LogInfo("read conf success path:conf/doc/%s")
	}
	return obj.(*Config%s)
}
'''%(structName,keyStr,structName,structName,tempfilename, structName, tempfilename, tempfilename, structName)
	return str




def GenConfig(indir,outdir):
	docs = []
	funDefInfo = {}
	fixedDefInfo = {}
	selfDefStr=""
	globalSetStr=""
	configStructStr = ""
	configStructReadStr = ""
	configStructGetStr = ""
	configFieldStr = ""
	readConfName=""
	files = ListDirFile(indir)
	for fname in files:
		if fname.find('ExportSetting.csv') != -1:
			selfDefStr, selfDefInfo = genSelfDefSetting(fname)
			break
	for fname in files:
		if fname.find('ExportSetting.csv') != -1:
			continue
		if fname.find('.csv') == -1:
			continue
		if fname.find('globalSet.csv') != -1:
			globalSetStr = genGlobalSetConfig(fname)
		else:
			str1,str2,str3,str4,str5 = genConfig(fname,selfDefInfo,fixedDefInfo)
			configStructStr+=str1
			configStructReadStr+=str2
			configStructGetStr+=str3
			configFieldStr+=str4
			readConfName+=str5

	fixedDefInfoStr=""
	for info in selfDefInfo.values():
		if info[0] ==True:
			selfDefStr+= info[2]
		if info[3] ==True:
			selfDefStr+= info[5]

	for value in fixedDefInfo.values():
		fixedDefInfoStr+= value[1]

	with open("%s/auto_config.go" %outdir,"w+") as f2:
		f2.write(
'''package frame
import (
	"antnet"
	"errors"
	"reflect"
)%s
type configDoc struct {
	ErrStr map[uint16]string `json:"errstr"`
	GlobalSet *ConfigGlobalSet `json:"global"`%s
}%s
func ReadConfigDoc() *configDoc {
	return &configDoc{
		errMap,
		ReadConfigGlobalSet(),%s
	}
}

'''%(fixedDefInfoStr+selfDefStr+globalSetStr+configStructStr+configStructReadStr,configFieldStr,configStructGetStr,readConfName)
			)
