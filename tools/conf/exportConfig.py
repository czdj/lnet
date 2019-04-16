#_*_coding:utf-8_*_
import os
import os.path
import sys
import ConfigParser
import platform
import genconfig

#配置数据
configParse =ConfigParser.ConfigParser()
configParse.read("exportConfig.txt")


#svn用户名和密码
svnUserName = configParse.get("exportConfig","svnUserName")
svnPassword = configParse.get("exportConfig","svnPassword")
svnConfigPath = configParse.get("exportConfig","svnConfigPath")
autoConfigPath = configParse.get("exportConfig","autoConfigPath")
configPath = configParse.get("exportConfig","configPath")
autoGoFilePath = configParse.get("exportConfig","autoGoFilePath")

def IsWindows():
	sysType = platform.system()
	if sysType == "Windows":
		return True
	return False
	
def ListDirFile(rootdir):
	files = []
	list = os.listdir(rootdir) #列出文件夹下所有的目录与文件
	for i in range(0,len(list)):
		path = os.path.join(rootdir,list[i])
		if os.path.isfile(path):
			files.append(path)
	return files

#删除目录
def removeDir(path):
	rdCmd = ""
	if IsWindows() == False:
		rdCmd = "rm -rf "+path
	else:
		rdCmd = "rd /s/q " +path
	os.system(rdCmd)

def IsWindows():
	sysType = platform.system()
	if sysType == "Windows":
		return True
	return False

def CopyFile(sourcePath,targetPath):
	cpCmd = ""
	if IsWindows() == False:
		cpCmd = "cp -r "+sourcePath+" " +targetPath
	else:
		cpCmd = "XCOPY " +sourcePath +" " +targetPath + " /E/Y "
	os.system(cpCmd)

def MoveFile(sourcePath,targetPath):
	mvCmd = ""
	if IsWindows() == False:
		mvCmd = "mv "+sourcePath+" " +targetPath
	else:
		mvCmd = "MOVE " +sourcePath +" " +targetPath
	os.system(mvCmd)

#切换目录
def CdPath(path):
	os.chdir(path)

def CopyDir(sourceDir, targetDir):
	if not os.path.exists(targetDir):
		os.makedirs(targetDir)
	for f in os.listdir(sourceDir):
		sourceF = os.path.join(sourceDir, f)
		targetF = os.path.join(targetDir, f)
		if os.path.isfile(sourceF):
			open(targetF, "wb+").write(open(sourceF, "rb").read())
		if os.path.isdir(sourceF):
			CopyDir(sourceF, targetF)
			
			
#svn更新
def SvnUp(targetPath):
	svnUpCmd = "svn up "+ targetPath +" --username "+svnUserName+" --password "+svnPassword+" --force"
	os.system(svnUpCmd)
		
		
def main():
	SvnUp(svnConfigPath)
	print("please wait a moment to gen csv file.")
	removeDir(autoConfigPath)
	exportXlsxCmd = "dotnet ./ExportXlsx/ExportXlsx.dll --optionSetting=./ServerExportXlsxSetting.json"
	os.system(exportXlsxCmd)
	print("please wait a moment to copy csv to doc dir")
	CopyDir(autoConfigPath,configPath)
	print("start to gen auto_config.go ")
	genconfig.GenConfig(configPath,autoGoFilePath)
	os.system("pause")


if __name__ == "__main__":
	main()
	sys.exit(0)

