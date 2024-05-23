v_pause="true"
v_pause_branchs="master ad164"
curr_branch="master"

if [ "$v_pause" = "true" ];then
  echo "$v_pause_branchs" | while read -r word; do
    if [ "$word" = "$curr_branch" ];then
      echo "@ENV_r_rs=$curr_branch 已暂停新的触发"
      exit 1
    fi
  done
fi


# 定义一个包含以空格分隔的字符串
str="one two three four five"

OLD_IFS="$IFS"  #保存当前shell默认的分割符，一会要恢复回去
IFS=" "                  #将shell的分割符号改为，“”
words=($st)     #分割符是“，”，"hello,shell,split,test" 赋值给array 就成了数组赋值
IFS="$OLD_IFS"  #恢复shell默认分割符配置


# 遍历数组中的每个单词
for word in "${words[@]}"; do
  echo "Word: $word"
done