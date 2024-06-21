# java
参考:
- [Snailclimb/JavaGuide](https://github.com/Snailclimb/JavaGuide)

## insall
JRE和JDK的区别： JRE是Java Runtime Envrionment，是用来运行Java环境的，并不能用来开发；JDK是Java Development Kit，是Java的开发组件，既能运行Java程序又能进行开发
JRE中带不带headless的区别： 带headless的是用来运行不包含GUI的java程序的，不带headless的可以运行带GUI的java程序

### repo源安装 
```bash
# apt install openjdk-8-jre/openjdk-8-jdk
# sudo update-alternatives --config java # 切换java version
```
### 使用ppa
1. 使用ppa/源方式安装
```
sudo add-apt-repository ppa:webupd8team/java
sudo apt-get update
```
### 安装oracle-java-installer
参考:
  - [Ubuntu 安装 JDK 7 / JDK8 的两种方式](http://www.cnblogs.com/a2211009/p/4265225.html)

```
# jdk7
sudo apt-get install oracle-java7-installer
#　jdk8
sudo apt-get install oracle-java8-installer
```

  如果因为防火墙或者其他原因,导致相应的installer下载速度很慢,可以中断操作，然后先下载好相应jdk的tar.gz包,放在文件夹:

  /var/cache/oracle-jdk7-installer             (jdk7)

  /var/cache/oracle-jdk8-installer             (jdk8)

  下面,再安装一次installer即可．此时installer会默认使用该tar.gz包．

1. 设置系统默认jdk
```
# jdk7
sudo update-java-alternatives -s java-7-oracle
# jdk8
sudo update-java-alternatives -s java-8-oracle
```

1. 测试jdk 是是否安装成功
```
java -version
javac -version
```

## 前提
server环境jdk可用: [dragonwell](https://dragonwell-jdk.io/)

`javap`是jdk自带的反解析工具。它的作用就是根据class字节码文件,反解析出当前类对应的code区(汇编指令)、本地变量表、异常表和代码行偏移量映射表、常量池等等信息.


一个java文件可以由多个class定义, 但至多有一个`public class`类.

jshell是 Java 9 新增的一个交互式的编程环境工具。它允许你无需使用类或者方法包装来执行 Java 语句, 可用于验证代码, 但实际可用性不高.

CLASSPATH:是由JRE提供的，用于定义Java程序解释时类加载路径，默认设置的是当前所在目录.

## base
java只有两种类型: 基本类型和class.

和其他编程语言相同, `main()是入口`.

已分号结尾的一行或多行代码被成为java语句.

java语法允许将低精度的值赋给高精度的值.

### 基本类型
Java 中有 8 种基本数据类型, 分别为:
- 6 种数字类型 :byte、short、int、long、float、double
- 1 种字符类型:char
- 1 种布尔型:boolean

这 8 种基本数据类型的默认值以及所占空间的大小如下:

基本类型    位数  字节  默认值
int 32  4   0
short   16  2   0
long    64  8   0L
byte    8   1   0
char    16  2   'u0000'
float   32  4   0f
double  64  8   0d
boolean 1       false

> 另外, 对于 boolean, 官方文档未明确定义, 它依赖于 JVM 厂商的具体实现. 逻辑上理解是占用 1 位, 但是实际中会考虑计算机高效存储因素.

基本类型用`===`判断相等.

这八种基本类型都有对应的包装类分别为:Byte、Short、Integer、Long、Float、Double、Character、Boolean.

包装类型不赋值就是 Null , 而基本类型有默认值且不是 Null. 另外, 这个问题建议还可以先从 JVM 层面来分析: 基本数据类型直接存放在 Java 虚拟机栈中的局部变量表中, 而包装类型属于对象类型, 而对象实例都存在于堆中. 相比于对象类型,  基本数据类型占用的空间非常小.

Java 基本类型的包装类的大部分都实现了常量池技术. Byte,Short,Integer,Long 这 4 种包装类默认创建了数值 [-128, 127] 的相应类型的缓存数据, Character 创建了数值在[0,127]范围的缓存数据, Boolean 直接返回 True Or False. 如果超出对应范围仍然会去创建新的对象, 缓存的范围区间的大小只是在性能和资源之间的权衡.

所有整型包装类对象之间值的比较, 全部使用 equals 方法比较, 避免常量池技术的干扰.

注意:
1. Java 里使用 long 类型的数据一定要在数值后面加上`L`, 否则将作为整型解析
1. char, 单引号包裹; String, 双引号包裹

### 自动装箱与拆箱
装箱:将基本类型用它们对应的引用类型包装起来
拆箱:将包装类型转换为基本数据类型

> 从字节码中, 我们发现装箱其实就是调用了 包装类的valueOf()方法, 拆箱其实就是调用了 xxxValue()方法

## class和object
> java中类是引用类型.

class是一种抽象和定义.
class的实例就是object.

new即是创建对象并返回引用, 对象在heap上.

重载是在类中重用方法名, 方法签名不同, 与返回值类型不相关, 但**推荐返回相同的类型.**

StringBuffer用于可变长度的字符串.

object常用方法:
- hashcode() : 根据对象所在的内存地址返回一个能够唯一代表这个对象的int值
- equals() : 判断同一个类的两个对象是否相等, 即判断它们是否指向同一对象.

### 继承
基本概念
1. java使用单继承
1. 子类继承了其父类中不是私有的成员变量和成员方法，作为自己的成员变量和方
1. 子类中定义的成员变量和父类中定义的成员变量相同时，则父类中的成员变量不能被继承
1. 子类中定义的成员方法，并且这个方法的名字返回类型，以及参数个数和类型与父类的某个成员方法完全相同，则父类的成员方法不能被继承

Object类是java中所有类的基类， 其他类都是直接或间接继承自它.

java创建一个子类对象s时会同时创建一个父类对象p, 可认为p内嵌在子类中. 通过子类的引用使用到其父类的属性时, 即访问了这个内嵌的父类对象的属性.

子类对象s赋值给父类对象p:
1. 编译时声明的对象是父类对象p, 但运行时却是使用子类对象s
1. p调用方法: 如果子类没有重写父类的法子是调用父类方法; 否则使用子类方法. 这两种情况下的this均是s
1. p使用属性: 无论子/父类是否存在相同属性, p均使用父类属性

> 需注意: 向上造型

### 多态
多态本质: 同一个东西在不同的环境下有不同表现. 多态在java的表现形式是重载和覆盖.

通过instanceof可判断一个引用指向的对象所属的类是不是某个类或某个类的子类.

### 方法
方法可见性控制:
1. public : 没有限制
1. 默认 : 只能被同一个包中的类使用
1. protected : 只能被子类使用
1. private: 只能被本类使用

### 内部类
内部类是在成员内部的类, 定义在类的内部是成员内部类, 定义在方法中的内部类是局部内部类.

## 泛型
泛型的本质是参数化类型, 也就是说所操作的数据类型被指定为一个参数.

Java 的泛型是伪泛型, 这是因为 Java 在编译期间, 所有的泛型信息都会被擦掉, 这也就是通常所说类型擦除.

泛型一般有三种使用方式:
1. 泛型类

    ```java
    //此处T可以随便写为任意标识, 常见的如T、E、K、V等形式的参数常用于表示泛型
    //在实例化泛型类时, 必须指定T的具体类型
    public class Generic<T>{

        private T key;

        public Generic(T key) {
            this.key = key;
        }

        public T getKey(){
            return key;
        }
    }

    Generic<Integer> genericInteger = new Generic<Integer>(123456);
    ```
1. 泛型接口

    ```java
    public interface Generator<T> {
        public T method();
    }

    // 实现泛型接口, 不指定类型
    class GeneratorImpl<T> implements Generator<T>{
        @Override
        public T method() {
            return null;
        }
    }
    // 实现泛型接口, 指定类型
    class GeneratorImpl<T> implements Generator<String>{
        @Override
        public String method() {
            return "hello";
        }
    }
    ```
1. 泛型方法


    ```java
    public static <E> void printArray( E[] inputArray )
   {
         for ( E element : inputArray ){
            System.out.printf( "%s ", element );
         }
         System.out.println();
    }

    // 创建不同类型数组: Integer, Double 和 Character
    Integer[] intArray = { 1, 2, 3 };
    String[] stringArray = { "Hello", "World" };
    printArray( intArray  );
    printArray( stringArray  );
    ```

## 反射
反射之所以被称为框架(Spring/Spring Boot、MyBatis等)的灵魂, 主要是因为它赋予了开发者在运行时分析类以及执行类中方法的能力.

通过反射可以获取任意一个类的所有属性和方法, 还可以调用这些方法和属性.

反射机制优缺点
优点 : 可以让咱们的代码更加灵活、为各种框架提供开箱即用的功能提供了便利
缺点 :在运行时有了分析操作类的能力, 这同样也增加了安全问题. 比如可以无视泛型参数的安全检查(泛型参数的安全检查发生在编译时). 另外, 反射的性能也要稍差点, 不过, 对于框架来说实际是影响不大的

另外, 像 Java 中的一大利器 注解 的实现也用到了反射.

## 异常
在 Java 中, 所有的异常都有一个共同的祖先 java.lang 包中的 Throwable 类. Throwable 类有两个重要的子类 Exception(异常)和 Error(错误). Exception 能被程序本身处理(try-catch),  Error 是无法处理的(只能尽量避免).

Exception 和 Error 二者都是 Java 异常处理的重要子类, 各自都包含大量子类:

- Exception :程序本身可以处理的异常, 可以通过 catch 来进行捕获

    Exception 又可以分为 受检查异常(必须处理) 和 不受检查异常(可以不处理):
    1. Java 代码在编译过程中, 如果受检查异常没有被 catch/throw 处理的话, 就没办法通过编译.

        除了RuntimeException及其子类以外, 其他的Exception类及其子类都属于受检查异常 . 常见的受检查异常有: IO 相关的异常、ClassNotFoundException 、SQLException等
    1. Java 代码在编译过程中 , 即使不处理不受检查异常也可以正常通过编译

        RuntimeException 及其子类都统称为非受检查异常, 例如:NullPointerException、NumberFormatException(字符串转换为数字)、ArrayIndexOutOfBoundsException(数组越界)、ClassCastException(类型转换错误)、ArithmeticException(算术错误)等
- Error :Error 属于程序无法处理的错误 , 没办法通过 catch 来进行捕获

    例如, Java 虚拟机运行错误(Virtual MachineError)、虚拟机内存不够错误(OutOfMemoryError)、类定义错误(NoClassDefFoundError)等 . 这些异常发生时, Java 虚拟机(JVM)一般会选择线程终止. 

Throwable 类常用方法:
- public string getMessage():返回异常发生时的简要描述
- public string toString():返回异常发生时的详细信息
- public string getLocalizedMessage():返回异常对象的本地化信息. 使用 Throwable 的子类覆盖这个方法, 可以生成本地化信息. 如果子类没有覆盖该方法, 则该方法返回的信息与 getMessage()返回的结果相同
- public void printStackTrace():在控制台上打印 Throwable 对象封装的异常信息

try-catch-finally:
- try块: 用于捕获异常. 其后可接零个或多个 catch 块, 如果没有 catch 块, 则必须跟一个 finally 块.
- catch块: 用于处理 try 捕获到的异常.
- finally 块: 无论是否捕获或处理异常, finally 块里的语句都会被执行. 当在 try 块或 catch 块中遇到 return 语句时, finally 语句块将在方法返回之前被执行.

在以下 3 种特殊情况下, finally 块不会被执行:
1. 在 try 或 finally块中用了 System.exit(int)退出程序. 但是, 如果 System.exit(int) 在异常语句之后, finally 还是会被执行
1. 程序所在的线程死亡
1. 关闭 CPU

注意: 当 try 语句和 finally 语句中都有 return 语句时, 在方法返回之前, finally 语句的内容将被执行, 并且 finally 语句的返回值将会覆盖原始的返回值

使用 try-with-resources 来代替try-catch-finally
- 适用范围(资源的定义): 任何实现 java.lang.AutoCloseable或者 java.io.Closeable 的对象
- 关闭资源和 finally 块的执行顺序: 在 try-with-resources 语句中, 任何 catch 或 finally 块在声明的资源关闭后运行

> 当然多个资源需要关闭的时候, 使用 try-with-resources 实现起来也非常简单, 如果还是用try-catch-finally可能会带来很多问题

## I/O
Java 中 IO 流分为几种?
- 按照流的流向分, 可以分为输入流和输出流
- 按照操作单元划分, 可以划分为字节流和字符流
- 按照流的角色划分为节点流和处理流

Java Io 流共涉及 40 多个类, 这些类看上去很杂乱, 但实际上很有规则, 而且彼此之间存在非常紧密的联系,  Java I0 流的 40 多个类都是从如下 4 个抽象类基类中派生出来的:
- InputStream/Reader: 所有的输入流的基类, 前者是字节输入流, 后者是字符输入流
- OutputStream/Writer: 所有输出流的基类, 前者是字节输出流, 后者是字符输出流

## FAQ
### JVM vs JDK vs JRE
Java 虚拟机(JVM)是运行 Java 字节码的虚拟机. JVM 有针对不同系统的特定实现(Windows, Linux, macOS), 目的是使用相同的字节码, 它们都会给出相同的结果. 字节码和不同系统的 JVM 实现是 Java 语言"一次编译, 随处可以运行"的关键所在.

> 在 Java 中, JVM 可以理解的代码就叫做字节码(即扩展名为 .class 的文件), 它不面向任何特定的处理器, 只面向虚拟机.

> JIT 编译器, 而 JIT 属于运行时编译. 当 JIT 编译器完成第一次编译后, 其会将字节码对应的机器码保存下来, 下次可以直接使用.

> JDK 9 引入了一种新的编译模式 AOT(Ahead of Time Compilation), 它是直接将字节码编译成机器码, 这样就避免了 JIT 预热等各方面的开销.

JDK 是 Java Development Kit 缩写, 它是功能齐全的 Java SDK. 它拥有 JRE 所拥有的一切, 还有编译器(javac)和工具(如 javadoc 和 jdb). 它能够创建和编译程序.

JRE 是 Java 运行时环境. 它是运行已编译 Java 程序所需的所有内容的集合, 包括 Java 虚拟机(JVM), Java 类库, java 命令和其他的一些基础构件. 但是, 它不能用于创建新程序.

### java字符型常量和字符串常量的区别?
形式 : 字符常量是单引号引起的一个字符, 字符串常量是双引号引起的 0 个或若干个字符
含义 : 字符常量相当于一个整型值( ASCII 值),可以参加表达式运算; 字符串常量代表一个地址值(该字符串在内存中存放位置)
占内存大小 : 字符常量只占 2 个字节; 字符串常量占若干个字节 (注意:char 在 Java 中占两个字节)

### Java 和 C++的区别?
都是面向对象的语言, 都支持封装、继承和多态

Java 不提供指针来直接访问内存, 程序内存更加安全
Java 的类是单继承的, C++ 支持多重继承；虽然 Java 的类不可以多继承, 但是接口可以多继承
Java 有自动内存管理垃圾回收机制(GC), 不需要程序员手动释放无用内存.
C ++同时支持方法重载和操作符重载, 但是 Java 只支持方法重载(操作符重载增加了复杂性, 这与 Java 最初的设计思想不符)

### ==和 equals 的区别
对于基本数据类型来说, ==比较的是值. 对于引用数据类型来说, ==比较的是对象的内存地址. 它们本质比较的都是值, 只是引用类型变量存的值是对象的地址.

equals() 作用不能用于判断基本数据类型的变量, 只能用来判断两个对象是否相等. equals()方法存在于Object类中, 而Object类是所有类的直接或间接父类.

equals() 方法存在两种使用情况:
- 类没有覆盖 equals()方法 :通过equals()比较该类的两个对象时, 等价于通过“==”比较这两个对象, 使用的默认是 Object类equals()方法
- 类覆盖了 equals()方法 :一般我们都覆盖 equals()方法来比较两个对象中的属性是否相等；若它们的属性相等, 则返回 true(即, 认为这两个对象相等)

### hashCode()与 equals(), 为什么重写 equals 时必须重写 hashCode 方法?
hashCode() 的作用是获取哈希码, 也称为散列码；它实际上是返回一个 int 整数. 这个哈希码的作用是确定该对象在哈希表中的索引位置. hashCode()定义在 JDK 的 Object 类中, 这就意味着 Java 中的任何类都包含有 hashCode() 函数. 另外需要注意的是: Object 的 hashcode 方法是本地方法, 也就是用 c 语言或 c++ 实现的, 该方法通常用来将对象的 内存地址 转换为整数之后返回.

如果两个对象相等, 则 hashcode 一定也是相同的. 两个对象相等,对两个对象分别调用 equals 方法都返回 true. 但是, 两个对象有相同的 hashcode 值, 它们也不一定是相等的(hash碰撞) . 因此, equals 方法被覆盖过, 则 hashCode 方法也必须被覆盖.

### 在一个静态方法内调用一个非静态成员为什么是非法的?
静态方法是属于类的, 在类加载的时候就会分配内存, 可以通过类名直接访问. 而非静态成员属于实例对象, 只有在对象实例化之后才存在, 然后通过类的实例对象去访问. 在类的非静态成员不存在的时候静态成员就已经存在了, 此时调用在内存中还不存在的非静态成员, 属于非法操作.

### 静态方法和实例方法有何不同?
1. 在外部调用静态方法时, 可以使用"类名.方法名"的方式, 也可以使用"对象名.方法名"的方式. 而实例方法只有后面这种方式. 也就是说, 调用静态方法可以无需创建对象
1. 静态方法在访问本类的成员时, 只允许访问静态成员(即静态成员变量和静态方法), 而不允许访问实例成员变量和实例方法；实例方法则无此限制

**java的方法必须是属于一个类或者类的实例的**

### 为什么 Java 中只有值传递?
按值调用(call by value) 表示方法接收的是调用者提供的值, 按引用调用(call by reference) 表示方法接收的是调用者提供的变量地址. 一个方法可以修改传递引用所对应的变量值, 而不能修改传递值调用所对应的变量值. 它用来描述各种程序设计语言(不只是 Java)中方法参数传递方式.

Java 程序设计语言总是采用按值调用. 也就是说, 方法得到的是所有参数值的一个拷贝, 也就是说, 方法不能修改传递给它的任何参数变量的内容.

### 重载和重写的区别
重载就是同一个类中多个同名方法根据不同的传参来执行不同的逻辑处理

重写就是子类对父类方法的重新改造, 外部样子不能改变, 内部逻辑可以改变.

区别点 重载方法    重写方法
- 发生范围    同一个类    子类
- 参数列表    必须修改    一定不能修改
- 返回类型    可修改 子类方法返回值类型应比父类方法返回值类型更小或相等
- 异常  可修改 子类方法声明抛出的异常类应比父类方法声明抛出的异常类更小或相等；
- 访问修饰符   可修改 一定不能做更严格的限制(可以降低限制)
- 发生阶段    编译期 运行期

方法的重写要遵循"两同两小一大":
- “两同”即方法名相同、形参列表相同；
- “两小”指的是子类方法返回值类型应比父类方法返回值类型更小或相等, 子类方法声明抛出的异常类应比父类方法声明抛出的异常类更小或相等
- “一大”指的是子类方法的访问权限应比父类方法的访问权限更大或相等

### 深拷贝 vs 浅拷贝
浅拷贝:对基本数据类型进行值传递, 对引用数据类型进行引用传递般的拷贝, 此为浅拷贝
深拷贝:对基本数据类型进行值传递, 对引用数据类型, 创建一个新的对象, 并复制其内容, 此为深拷贝

### 面向对象和面向过程的区别
- 面向过程 :**通常**面向过程性能比面向对象高

    因为类调用时需要实例化, 开销比较大, 比较消耗资源, 所以当性能是最重要的考量因素的时候, 比如单片机、嵌入式开发、Linux/Unix 等一般采用面向过程开发. 但是, 面向过程没有面向对象易维护、易复用、易扩展.
- 面向对象 :面向对象易维护、易复用、易扩展

    因为面向对象有封装、继承、多态性的特性, 所以可以设计出低耦合的系统, 使系统更加灵活、更加易于维护. 但是, 面向对象性能比面向过程低

### 成员变量与局部变量的区别有哪些?
1. 从语法形式上看, 成员变量是属于类的, 而局部变量是在代码块或方法中定义的变量或是方法的参数；成员变量可以被 public,private,static 等修饰符所修饰, 而局部变量不能被访问控制修饰符及 static 所修饰；但是, 成员变量和局部变量都能被 final 所修饰
1. 从变量在内存中的存储方式来看,如果成员变量是使用 static 修饰的, 那么这个成员变量是属于类的, 如果没有使用 static 修饰, 这个成员变量是属于实例的. 而对象存在于堆内存, 局部变量则存在于栈内存.
1. 从变量在内存中的生存时间上看, 成员变量是对象的一部分, 它随着对象的创建而存在, 而局部变量随着方法的调用而自动消失
1. 从变量是否有默认值来看, 成员变量如果没有被赋初, 则会自动以类型的默认值而赋值(一种情况例外:被 final 修饰的成员变量也必须显式地赋值), 而局部变量则不会自动赋值

### 创建一个对象用什么运算符?对象实体与对象引用有何不同?
new 运算符, new 创建对象实例(对象实例在堆内存中), 对象引用指向对象实例(对象引用存放在栈内存中).

### 一个类的构造方法的作用是什么? 若一个类没有声明构造方法, 该程序能正确执行吗? 为什么?
构造方法主要作用是完成对类对象的初始化工作.

如果一个类没有声明构造方法, 也可以执行！因为一个类即使没有声明构造方法也会有默认的不带参数的构造方法. 如果开发者添加了类的构造方法(无论是否有参), Java 就不会再添加默认的无参数的构造方法了, 这时候, 就不能直接 new 一个对象而不传递参数了, 所以我们一直在不知不觉地使用构造方法, 这也是为什么我们在创建对象的时候后面要加一个括号(因为要调用无参的构造方法). **如果我们重载了有参的构造方法, 记得都要把无参的构造方法也写出来(无论是否用到), 因为这可以帮助我们在创建对象的时候少踩坑**.

### 构造方法有哪些特点?是否可被 override?
特点:
- 名字与类名相同
- 没有返回值, 但不能用 void 声明构造函数
- 生成类的对象时自动执行, 无需调用

构造方法不能被 override(重写),但是可以 overload(重载),所以可以看到一个类中有多个构造函数的情况

### 面向对象三大特征
1. 封装

    把一个对象的状态信息(也就是属性)隐藏在对象内部, 不允许外部对象直接访问对象的内部信息. 但是可以提供一些可以被外界访问的方法来操作属性. 
1. 继承

    不同类型的对象, 相互之间经常有一定数量的共同点. 同时, 每一个对象还定义了额外的特性使得他们与众不同. 继承是使用已存在的类的定义作为基础建立新类的技术, 新类的定义可以增加新的数据或新的功能, 也可以用父类的功能, 但不能选择性地继承父类. 通过使用继承, 可以快速地创建新的类, 可以提高代码的重用, 程序的可维护性, 节省大量创建新类的时间 , 提高我们的开发效率.

    关于继承如下 3 点请记住:

    1. 子类拥有父类对象所有的属性和方法(包括私有属性和私有方法), 但是父类中的私有属性和方法子类是无法访问, 只是拥有
    1. 子类可以拥有自己属性和方法, 即子类可以对父类进行扩展
    1. 子类可以用自己的方式实现父类的方法
1. 多态

    多态, 顾名思义, 表示一个对象具有多种的状态. **具体表现为父类的引用指向子类的实例**

    它在继承的基础上扩充而来, 指的是类型的转换处理.

    多态的特点:

    1. 对象类型和引用类型之间具有继承(类)/实现(接口)的关系
    1. 引用类型变量发出的方法调用的到底是哪个类中的方法, 必须在程序运行期间才能确定
    1. 多态不能调用“只在子类存在但在父类不存在”的方法
    1. 如果子类重写了父类的方法, 真正执行的是子类覆盖的方法, 如果子类没有覆盖父类的方法, 执行的是父类的方法

### String StringBuffer 和 StringBuilder 的区别是什么? String 为什么是不可变的?
在 Java 9 之后, String 、StringBuilder 与 StringBuffer 的实现改用 byte 数组存储字符串. String 类中使用 final 关键字修饰用来保存字符串的数组, 所以String 对象是不可变的.

而 StringBuilder 与 StringBuffer 都继承自 AbstractStringBuilder 类, 在 AbstractStringBuilder 中没有用 final 关键字修饰保存字符串的数组, 所以这两种对象都是可变的.

String 中的对象是不可变的, 也就可以理解为常量, 线程安全. StringBuffer 对方法加了同步锁或者对调用的方法加了同步锁, 所以是线程安全的. StringBuilder 并没有对方法进行加同步锁, 所以是非线程安全的.

每次对 String 类型进行改变的时候, 都会生成一个新的 String 对象, 然后将指针指向新的 String 对象. StringBuffer 每次都会对 StringBuffer 对象本身进行操作, 而不是生成新的对象并改变对象引用. 相同情况下使用 StringBuilder 相比使用 StringBuffer 仅能获得 10%~15% 左右的性能提升, 但却要冒多线程不安全的风险.

对于三者使用的总结:
1. 操作少量的数据: 适用 String
1. 单线程操作字符串缓冲区下操作大量数据: 适用 StringBuilder
1. 多线程操作字符串缓冲区下操作大量数据: 适用 StringBuffer

### Java 序列化中如果有些字段不想进行序列化, 怎么办?
对于不想进行序列化的变量, 使用transient关键字修饰.

transient 关键字的作用是:阻止实例中那些用此关键字修饰的的变量序列化；当对象被反序列化时, 被 transient 修饰的变量值不会被持久化和恢复. transient 只能修饰变量, 不能修饰类和方法.

### 既然有了字节流,为什么还要有字符流?
问题本质想问:不管是文件读写还是网络发送接收, 信息的最小存储单元都是字节, 那为什么 I/O 流操作要分为字节流操作和字符流操作呢?

回答:字符流是由 Java 虚拟机将字节转换得到的, 问题就出在这个过程还算是非常耗时, 并且, 如果不知道编码类型就很容易出现乱码问题. 所以,  I/O 流就干脆提供了一个直接操作字符的接口, 方便我们平时对字符进行流操作. 如果音频文件、图片等媒体文件用字节流比较好, 如果涉及到字符的话使用字符流比较好.

### this/super
this关键字用于引用类的当前实例, 即指向本类对象.

super关键字用于从子类访问父类的变量和方法, 即指向父类对象.

### final/static
final关键字, 意思是最终的、不可修改的, 最见不得变化 , 用来修饰类、方法和变量, 具有以下特点:
- final修饰的类不能被继承, final类中的所有成员方法都会被隐式的指定为final方法
- final修饰的方法不能被重写
- final修饰的变量是常量, 如果是基本数据类型的变量, 则其数值一旦在初始化之后便不能更改；如果是引用类型的变量, 则在对其初始化之后便不能让其指向另一个对象

> 类中所有的private方法都隐式地指定为了final.

static 关键字主要有以下四种使用场景:
1. 修饰成员变量和成员方法: 被 static 修饰的成员属于类, 不属于单个这个类的某个对象, 被**该类的所有对象共享**, 可以并且建议通过类名调用。被static 声明的成员变量属于静态成员变量(static variable), 静态变量 存放在 Java 内存区域的方法区; 使用static修饰的方法是静态方法(static method)。调用格式:`类名.静态变量名 类名.静态方法名()`

    静态变量在程序加载时创建, 先于类的所有对象的创建.

    静态方法不能使用非静态的变量和方法, 因此是不能使用this, 因为this代表对象, 非静态. 静态方法也不存在覆盖.
1. 静态代码块: 静态代码块定义在类中方法外, 静态代码块在非静态代码块之前执行(静态代码块—>非静态代码块—>构造方法). 该类不管创建多少对象, 静态代码块只执行一次.
1. 静态内部类(static修饰类的话只能修饰内部类): 静态内部类与非静态内部类之间存在一个最大的区别: 非静态内部类在编译完成之后会隐含地保存着一个引用, 该引用是指向创建它的外围类, 但是静态内部类却没有。没有这个引用就意味着:1. 它的创建是不需要依赖外围类的创建。2. 它不能使用任何外围类的非static成员变量和方法
1. 静态导包(用来导入类中的静态资源, 1.5之后的新特性): 格式为:import static 这两个关键字连用可以指定导入某个类中的指定静态资源, 并且不需要使用类名调用类中静态成员, 可以直接使用类中静态成员变量和成员方法

### Interface、extends、implement的区别
interface是定义接口的关键字
implement是实现接口的关键字
extends是子类继承父类的关键字

interface中的变量必须使用static和final修饰. 接口中的方法只能是抽象方法, 必须使用public和abstract修饰.

### transient属性
一个对象只要实现了Serilizable接口，这个对象就可以被序列化， 此时将不需要序列化的属性前添加关键字transient，序列化对象的时候，这个属性就不会序列化到指定的目的地中.

### implements Serializable, Cloneable
Cloneable接口与Serializable接口都是定义接口而没有任何的方法. Cloneable可以实现对象的克隆复制，Serializable主要是对象序列化的接口定义. 很多时候我们涉及到对象的复制, 我们不可能都去使用setter去实现，这样编写代码的效率太低. JDK提供的Cloneable接口正是为了解决对象复制的问题而存在. Cloneable结合Serializable接口可以实现JVM对象的深度复制.

Cloneable接口是一个空接口，仅用于标记对象，Cloneable接口里面是没有clone()方法，的clone()方法是Object类里面的方法！默认实现是一个Native方法
```java
protected native Object clone() throws CloneNotSupportedException;
```
如果对象implement Cloneable接口的话，需要覆盖clone方法（因为Object类的clone方法是protected，需要覆盖为public）
```java
public Object clone() throws CloneNotSupportedException{
    return super.clone();
}
```
Object类里的clone()方法仅仅用于**浅拷贝**（拷贝基本成员属性，对于引用类型仅返回指向改地址的引用.


### [static](https://zhuanlan.zhihu.com/p/70110497)
静态变量存放在java vm的方法区中，并且是被所有线程所共享的.

static关键字总结:

1. 特点:
　　1. static是一个修饰符，用于修饰成员。（成员变量，成员函数）static修饰的成员变量 称之为静态变量或类变量.
　　2. static修饰的成员被所有的对象共享.
　　3. static优先于对象存在，因为static的成员随着类的加载就已经存在.
　　4. static修饰的成员多了一种调用方式，可以直接被类名所调用，（类名.静态成员）.
　　5. static修饰的数据是共享数据，对象中的存储的是特有的数据.

1. 成员变量和静态变量的区别:
　　1. 生命周期的不同:

　　　　成员变量随着对象的创建而存在随着对象的回收而释放.

　　　　静态变量随着类的加载而存在随着类的消失而消失.
　　2. 调用方式不同:

　　　　成员变量只能被对象调用.

　　　　静态变量可以被对象调用，也可以用类名调用.（推荐用类名调用）
　　3. 别名不同:

　　　　成员变量也称为实例变量.

　　　　静态变量称为类变量.
　　4. 数据存储位置不同:

　　　　成员变量数据存储在堆内存的对象中，所以也叫对象的特有数据.

　　　　静态变量数据存储在方法区（共享数据区）的静态区，所以也叫对象的共享数据.

1. 静态使用时需要注意的事项:

　　1. 静态方法只能访问静态成员（非静态既可以访问静态，又可以访问非静态）
　　2. 静态方法中不可以使用this或者super关键字
　　3. java主函数是静态的

Java里静态语句块是优先对象存在，也就是优先于构造方法存在，我们通常用来做只创建一次对象使用，类似于单列模式而且执行的顺序是：父类静态语句块 -> 子类静态语句块 -> 父类构造方法 -> 子类构造方法

```bash
$ vim TestMethod.java

public class TestMethod extends BaseClass {
 
    static int a;
 
    public TestMethod() {
        super();
        System.out.println("constructor of exec");
    }
 
    static {
        String a1="12";
        String a2="22";
        a=Integer.parseInt(a1)+Integer.parseInt(a2);
        System.out.println("chilren static block");
    }
 
    public static void main(String[] args) {
        System.out.println(TestMethod.a);
        TestMethod.a=45;
        new TestMethod();
        System.out.println(TestMethod.a);
        new TestMethod();
        System.out.println(TestMethod.a);
    }
 
}
 
class BaseClass{
    
    static int a;
    
    static {
        String a1="10";
        String a2="20";
        a=Integer.parseInt(a1)+Integer.parseInt(a2);
        System.out.println("baseClass static block");
        System.out.println(a);
    }
    
    public BaseClass(){
        System.out.println("Base class constructor of exec");
        System.out.println(BaseClass.a);
        BaseClass.a=300;
        System.out.println(BaseClass.a);
    }
}
$ javac TestMethod.java
$ java TestMethod
baseClass static block
30
chilren static block
34
Base class constructor of exec # 开始new TestMethod()
30
300
constructor of exec
45
Base class constructor of exec # 再次new TestMethod(), 是用了同一个父类
300
300
constructor of exec
45
```

### java properties文件加载包含反斜杠
在java中，利用Properties.load()加载配置文件时，如果配置文件含有"\", 则会将反斜杠作为转义符处理，而不是作为正常字符.

### volatile
volatile是Java提供的一种轻量级的同步机制。Java 语言包含两种内在的同步机制：同步块（或方法）和 volatile 变量，相比于synchronized（synchronized通常称为重量级锁），volatile更轻量级，因为它不会引起线程上下文的切换和调度。但是volatile 变量的同步性较差（有时它更简单并且开销更低），而且其使用也更容易出错.

volatile可以保证线程可见性且提供了一定的有序性，但是无法保证原子性。在JVM底层volatile是采用“内存屏障”来实现的。观察加入volatile关键字和没有加入volatile关键字时所生成的汇编代码发现，加入volatile关键字时，会多出一个lock前缀指令，lock前缀指令实际上相当于一个内存屏障（也成内存栅栏），内存屏障会提供3个功能：

1. 它确保指令重排序时不会把其后面的指令排到内存屏障之前的位置，也不会把前面的指令排到内存屏障的后面；即在执行到内存屏障这句指令时，在它前面的操作已经全部完成；
1. 它会强制将对缓存的修改操作立即写入主存；
1. 如果是写操作，它会导致其他CPU中对应的缓存行无效

### jdbc
在Spring中，通过JDBC驱动定义数据源是最简单的配置方式。Spring提供了三个这样的数据源类供选择：
- DriverManagerDataSource：在每个连接请求时都会返回一个新建的连接。与JDBC的BasicDataSource不同，由DriverManagerDataSource提供的连接并没有进行池化管理。
- SimpleDriverDataSource：与DriverManagerDataSource工作方式类似，但是它直接使用JDBC驱动，来解决在特定环境下的类加载问题，这样的环境包括OSGi容器。
- SingleConnectionDataSource：在每个连接请求时都会返回同一个的连接。尽管SingleConnectionDataSource不是严格意义上的连接池数据源，但是可以将其视为只有一个连接的池。

注意：**SingleConnectionDataSource有且只有一个数据库连接，不适于多线程**，DriverManagerDataSource和SimpleDriverDataSource尽管支持多线程，但是在每次请求的时候都会创建新连接，这是以性能为代价的.

JdbcTemplate是Spring对JDBC的封装，目的是使JDBC更加易于使用。JdbcTemplate是Spring的一部分。JdbcTemplate处理了资源的建立和释放。他帮助我们避免一些常见的错误，比如忘了总要关闭连接。他运行核心的JDBC工作流，如Statement的建立和执行，而我们只需要提供SQL语句和提取结果.

### ExecutorService
ExecutorService是Java中对线程池定义的一个接口.

### @StaticMetamodel(Role.class)
参考:
- [Java EE 8 Tutorial by Oracle / Using the Metamodel API to Model Entity Classes](https://javaee.github.io/tutorial/persistence-criteria002.html)

这是JPA Criteria API的一部分

JPA Criteria API 是JPA中用于构建OOP风格的，较为灵活的SQL的一部分.其中包含Metamodel API，该API用于对Entity进行建模，JPA Criteria API有两种使用方式

强类型方式，即使用Metamodel来获得Entity的属性和域
基于字符串方式，直接使用属性名或域名来获得
@StaticMetamodel(Role.class)注解表明Role_是Role的Metamodel.也就是说在构建Criteria查询的时候如下的语句来指定where条件
```java
cq.where(role.get(Role_.roleName).in("admin", "superadmin"));
//相当于 ···WHERE role_name IN ('admin','superadmin')
//cq 是一个CriteriaQuery实例
```

### abstract
抽象方法：在类中没有方法体的方法，就是抽象方法，抽象方法用abstract修饰.

抽象类：在面向对象的概念中，如果一个类中没有包含足够的信息来描绘一个具体的对象，这样的类就是抽象类。

抽象类除了不能实例化对象之外，类的其它功能依然存在，成员变量、成员方法和构造方法的访问方式和普通类一样

因为无法执行抽象方法，因此这个类也必须申明为抽象类（abstract class）.

如果父类的方法本身不需要实现任何功能，仅仅是为了定义方法签名，目的是让子类去覆写它，那么，可以把父类的方法声明为抽象方法:
```java
abstract class Person {
    public abstract void run();
}
```

把一个方法声明为abstract，表示它是一个抽象方法，本身没有实现任何方法语句. 因为这个抽象方法本身是无法执行的，所以，Person类也无法被实例化. 编译器会告诉我们，无法编译Person类，因为它包含抽象方法, 因此必须把Person类本身也声明为abstract.

因为**抽象类本身被设计成只能用于被继承，因此，抽象类可以强迫子类实现其定义的抽象方法**，否则编译会报错. 因此，抽象方法实际上相当于定义了“规范”.

抽象类语法:
1. 抽象类是抽象的类型, 和接口一样, 在编程中不能通过new来创建其对象
1. 当一个类声明实现了某个或某些接口, 但却没有实现接口中定义的所有抽象方法时, 这个类就必须声明为抽象类

抽象类作用: 在语法层, 强制其子类覆盖某些方法.

### 抽象类和接口的区别
1. 语法层面上的区别
  1. 抽象类可以提供成员方法的实现细节，而接口中只能存在public abstract 方法；
  2. 抽象类中的成员变量可以是各种类型的，而接口中的成员变量只能是public static final类型的；
  3. 接口中不能含有静态代码块以及静态方法，而抽象类可以有静态代码块和静态方法；
  4. 一个类只能继承一个抽象类，而一个类却可以实现多个接口

2. 设计层面上的区别
  1. 抽象类是对一种事物的抽象，即对类抽象，而接口是对行为的抽象。抽象类是对整个类整体进行抽象，包括属性、行为，但是接口却是对类局部（行为）进行抽象。举个简单的例子，飞机和鸟是不同类的事物，但是它们都有一个共性，就是都会飞。那么在设计的时候，可以将飞机设计为一个类Airplane，将鸟设计为一个类Bird，但是不能将 飞行 这个特性也设计为类，因此它只是一个行为特性，并不是对一类事物的抽象描述。此时可以将 飞行 设计为一个接口Fly，包含方法fly( )，然后Airplane和Bird分别根据自己的需要实现Fly这个接口。然后至于有不同种类的飞机，比如战斗机、民用飞机等直接继承Airplane即可，对于鸟也是类似的，不同种类的鸟直接继承Bird类即可。从这里可以看出，继承是一个 "是不是"的关系，而 接口 实现则是 "有没有"的关系。如果一个类继承了某个抽象类，则子类必定是抽象类的种类，而接口实现则是有没有、具备不具备的关系，比如鸟是否能飞（或者是否具备飞行这个特点），能飞行则可以实现这个接口，不能飞行就不实现这个接口。

　2. 设计层面不同，抽象类作为很多子类的父类，它是一种模板式设计。而接口是一种行为规范，它是一种辐射式设计. 什么是模板式设计？最简单例子，大家都用过ppt里面的模板，如果用模板A设计了ppt B和ppt C，ppt B和ppt C公共的部分就是模板A了，如果它们的公共部分需要改动，则只需要改动模板A就可以了，不需要重新对ppt B和ppt C进行改动。而辐射式设计，比如某个电梯都装了某种报警器，一旦要更新报警器，就必须全部更新。也就是说对于抽象类，如果需要添加新的方法，可以直接在抽象类中添加具体的实现，子类可以不进行变更；而对于接口则不行，如果接口进行了变更，则所有实现这个接口的类都必须进行相应的改动。

### synchronized
synchronized是Java中的关键字，是一种同步锁, 可以保证方法或者代码块在运行时，同一时刻只有一个方法可以进入到临界区，同时它还可以保证共享变量的内存可见性。它修饰的对象有以下几种： 
1. 修饰一个代码块，被修饰的代码块称为同步语句块，其作用的范围是大括号{}括起来的代码，作用的对象是调用这个代码块的对象； 
2. 修饰一个方法，被修饰的方法称为同步方法，其作用的范围是整个方法，作用的对象是调用这个方法的对象； 
3. 修改一个静态的方法，其作用的范围是整个静态方法，作用的对象是这个类的所有对象； 
4. 修改一个类，其作用的范围是synchronized后面括号括起来的部分，作用主的对象是这个类的所有对象

类似go的Mutex.
### new interface
[反编译观察一下，发现原来是编译器自动生成一个类](https://www.cnblogs.com/yjmyzz/p/3448330.html).

### instanceof
测试它左边的对象是否是它右边的类的实例

### 子类赋值父类
super相当于指向父类示例的一个指针; 子类只保存子类的信息和super指针; 当JVM加载一个子类的时候也会把它的父类一同加载的，子类内部通过super保存父类的一个引用.

因此子类里有一个区域放的父类的实例，子类内存区里有一个指针，指向了这个内存区里包括的父类实例区，当把引用付给父类时，是把子类内存区里面的父类实例区域的引用给了父类的实例.

### threadlocal
threadlocal是一个线程内部的存储类,可以在指定线程内存储数据,数据存储以后,只有指定线程可以得到存储数据

## 多线程
继承java.lang.Thread, 并实现start(),run()即可.

方法:
- Thread.currentThread() : 获取当前线程

synchronized是方法的修饰符, 作用是每次只能有一个线程执行此类中的同步方法.

## 注解
### Java语言使用@interface语法来定义注解（Annotation）
注解：提供一种为程序元素设置元数据的方法.

基本原则：注解不能直接干扰程序代码的运行，无论增加或删除注解，代码都能够正常运行。
注解（也被成为元数据）为我们在代码中添加信息提供了一种形式化的方法，使我们可以在稍后某个时刻非常方便地使用这些数据。 ———摘自《Thinking in Java》

简单来说注解的作用就是将我们的需要的数据储存起来，在以后的某一个时刻（可能是编译时，也可能是运行时）去调用它.

```java
@Target(ElementType.TYPE) // 只能应用于类上
@Retention(RetentionPolicy.RUNTIME) // 保存到运行时
public @interface DBTable {
    String name() default "";
}

//在类上使用该注解
@DBTable(name = "MEMBER")
public class Member {
    //.......
}
```

上述定义一个名为DBTable的注解，该用于主要用于数据库表与Bean类的映射. 声明一个String类型的name元素，其默认值为空字符，但是必须注意到对应任何元素的声明应采用方法的声明方式，同时可选择使用default提供默认值.


#### 自定义注解
自定义注解，是使用元注解来实现的.

元注解:
- `@Target`说明了Annotation所修饰的对象范围

  Annotation可被用于 packages、types（类、接口、枚举、Annotation类型）、类型成员（方法、构造方法、成员变量、枚举值）、方法参数和本地变量（如循环变量、catch参数）. 在Annotation类型的声明中使用了target可更加明晰其修饰的目标.

  作用：用于描述注解的使用范围（即：被描述的注解可以用在什么地方）

  取值(ElementType)有：
  - CONSTRUCTOR:用于描述构造器
  - FIELD:用于描述域
  - LOCAL_VARIABLE:用于描述局部变量
  - METHOD:用于描述方法
  - PACKAGE:用于描述包
  - PARAMETER:用于描述参数
  - TYPE:用于描述类、接口(包括注解类型) 或enum声明
- `@Retention`定义了注解的保留策略即该Annotation被保留的时间长短

  某些Annotation仅出现在源代码中，而被编译器丢弃；而另一些却被编译在class文件中；编译在class文件中的Annotation可能会被虚拟机忽略，而另一些在class被装载时将被读取（请注意并不影响class的执行，因为Annotation与class在使用上是被分离的）。使用这个meta-Annotation可以对 Annotation的“生命周期”限制

  作用：表示需要在什么级别保存该注释信息，用于描述注解的生命周期（即：被描述的注解在什么范围内有效）

  取值（RetentionPoicy）有：
  - SOURCE:在源文件中有效（即源文件保留）
  - CLASS:在class文件中有效（即class保留）
  - RUNTIME:在运行时有效（即运行时保留）
- `@Documented`用于描述其它类型的annotation应该被作为被标注的程序成员的公共API，因此可以被例如javadoc此类的工具文档化。Documented是一个标记注解，没有成员
- `@Inherited`是一个标记注解，@Inherited阐述了某个被标注的类型是被继承的。如果一个使用了@Inherited修饰的annotation类型被用于一个class，则这个annotation将被用于该class的子类。

  注意：@Inherited annotation类型是被标注过的class的子类所继承。类并不从它所实现的接口继承annotation，方法并不从它所重载的方法继承annotation。

  当@Inherited annotation类型标注的annotation的Retention是RetentionPolicy.RUNTIME，则反射API增强了这种继承性。如果我们使用java.lang.reflect去查询一个@Inherited annotation类型的annotation时，反射代码检查将展开工作：检查class和其父类，直到发现指定的annotation类型被发现，或者到达类继承结构的顶层

使用@interface自定义注解时，自动继承了java.lang.annotation.Annotation接口，由编译程序自动完成其他细节。在定义注解时，不能继承其他的注解或接口。@interface用来声明一个注解，其中的每一个方法实际上是声明了一个配置参数。方法的名称就是参数的名称，返回值类型就是参数的类型（返回值类型只能是基本类型、Class、String、enum）。可以通过default来声明参数的默认值.

定义注解格式: `public @interface 注解名 {定义体}`

注解参数的可支持数据类型:
- 所有基本数据类型（int,float,boolean,byte,double,char,long,short)
- String类型
- Class类型
- enum类型
- Annotation类型
- 以上所有类型的数组

Annotation类型里面的参数该怎么设定:
- 首先,只能用public或默认(default)这两个访问权修饰.例如,String value();这里把方法设为defaul默认类型；　 　
- 其次,参数成员只能用基本类型byte,short,char,int,long,float,double,boolean八种基本数据类型和 String,Enum,Class,annotations等数据类型,以及这一些类型的数组.例如,String value();这里的参数成员就为String;　　
- 最后,如果只有一个参数成员,最好把参数名称设为”value”,后加小括号.例:下面的例子FruitName注解就只有一个参数成员。

### getSimpleName
```java
package com.example;

public final class TestClassNames {
    private static void showClass(Class<?> c) {
        System.out.println("getName():          " + c.getName());
        System.out.println("getCanonicalName(): " + c.getCanonicalName());
        System.out.println("getSimpleName():    " + c.getSimpleName());
        System.out.println("toString():         " + c.toString());
        System.out.println();
    }

    private static void x(Runnable r) {
        showClass(r.getClass());
        showClass(java.lang.reflect.Array.newInstance(r.getClass(), 1).getClass()); // Obtains an array class of a lambda base type.
    }

    public static class NestedClass {}

    public class InnerClass {}

    public static void main(String[] args) {
        class LocalClass {}
        showClass(void.class);
        showClass(int.class);
        showClass(String.class);
        showClass(Runnable.class);
        showClass(SomeEnum.class);
        showClass(SomeAnnotation.class);
        showClass(int[].class);
        showClass(String[].class);
        showClass(NestedClass.class);
        showClass(InnerClass.class);
        showClass(LocalClass.class);
        showClass(LocalClass[].class);
        Object anonymous = new java.io.Serializable() {};
        showClass(anonymous.getClass());
        showClass(java.lang.reflect.Array.newInstance(anonymous.getClass(), 1).getClass()); // Obtains an array class of an anonymous base type.
        x(() -> {});
    }
}

enum SomeEnum {
   BLUE, YELLOW, RED;
}

@interface SomeAnnotation {}
```

out:
```
getName():          void
getCanonicalName(): void
getSimpleName():    void
toString():         void

getName():          int
getCanonicalName(): int
getSimpleName():    int
toString():         int

getName():          java.lang.String
getCanonicalName(): java.lang.String
getSimpleName():    String
toString():         class java.lang.String

getName():          java.lang.Runnable
getCanonicalName(): java.lang.Runnable
getSimpleName():    Runnable
toString():         interface java.lang.Runnable

getName():          com.example.SomeEnum
getCanonicalName(): com.example.SomeEnum
getSimpleName():    SomeEnum
toString():         class com.example.SomeEnum

getName():          com.example.SomeAnnotation
getCanonicalName(): com.example.SomeAnnotation
getSimpleName():    SomeAnnotation
toString():         interface com.example.SomeAnnotation

getName():          [I
getCanonicalName(): int[]
getSimpleName():    int[]
toString():         class [I

getName():          [Ljava.lang.String;
getCanonicalName(): java.lang.String[]
getSimpleName():    String[]
toString():         class [Ljava.lang.String;

getName():          com.example.TestClassNames$NestedClass
getCanonicalName(): com.example.TestClassNames.NestedClass
getSimpleName():    NestedClass
toString():         class com.example.TestClassNames$NestedClass

getName():          com.example.TestClassNames$InnerClass // `$`表示嵌套的class
getCanonicalName(): com.example.TestClassNames.InnerClass
getSimpleName():    InnerClass
toString():         class com.example.TestClassNames$InnerClass

getName():          com.example.TestClassNames$1LocalClass
getCanonicalName(): null
getSimpleName():    LocalClass
toString():         class com.example.TestClassNames$1LocalClass

getName():          [Lcom.example.TestClassNames$1LocalClass;
getCanonicalName(): null
getSimpleName():    LocalClass[]
toString():         class [Lcom.example.TestClassNames$1LocalClass;

getName():          com.example.TestClassNames$1
getCanonicalName(): null
getSimpleName():    
toString():         class com.example.TestClassNames$1

getName():          [Lcom.example.TestClassNames$1;
getCanonicalName(): null
getSimpleName():    []
toString():         class [Lcom.example.TestClassNames$1;

getName():          com.example.TestClassNames$$Lambda$1/1175962212
getCanonicalName(): com.example.TestClassNames$$Lambda$1/1175962212
getSimpleName():    TestClassNames$$Lambda$1/1175962212
toString():         class com.example.TestClassNames$$Lambda$1/1175962212

getName():          [Lcom.example.TestClassNames$$Lambda$1;
getCanonicalName(): com.example.TestClassNames$$Lambda$1/1175962212[]
getSimpleName():    TestClassNames$$Lambda$1/1175962212[]
toString():         class [Lcom.example.TestClassNames$$Lambda$1;
```

# 构建
Apache Ant和Maven都是流行的构建工具，它们各自具有不同的优势和适用场景。 Ant的构建模型相对简单，适合较小的项目或需要简单构建过程的场景；而Maven的依赖管理功能和插件架构更加灵活和强大，更适合于大型项目或需要高度自动化的构建过程.

```bash
$ apt install ant
```

# java框架
## Spring
### ioc
生成流程: xml定义bean -> BeanDefinition(封装bean的定义信息) -> bean.

BeanDefinition->bean, 由BeanFactory生成, 创建bean实例由三种方法:
1. 反射
1. 工厂方法: @Bean
1. 工厂类: FactoryBean

bean实例创建后会利用`@Autowired`, `@Value`进行属性注入, 此时会利用三级缓存解决循环依赖.

之后bean会调用其生命周期的方法和aware, 可见BeanFactory注释.

> FactoryBean是一个特殊的接口，实现getObject()达到替换object的目的.

> @Autowired的原理: 在启动spring IoC时，容器自动装载了一个AutowiredAnnotationBeanPostProcessor后置处理器，当容器扫描到@Autowied、@Resource(是CommonAnnotationBeanPostProcessor后置处理器处理的)或@Inject时，就会在IoC容器自动查找需要的bean，并装配给该对象的属性.

### Spring MVC的web.xml配置详解
web.xml文件的作用是配置web工程启动,对于一个web工程来说，web.xml可以有也可以没有，如果存在web.xml文件；web工程在启动的时候，web容器(tomcat容器)会去加载web.xml文件，然后按照一定规则配置web.xml文件中的组件.


web容器加载顺序:ServletContext -> context-param(启动的初始化参数) -> listener -> filter ->servlet, **不会因在web.xml中的书写顺序改变**:
1. web容器启动后,会去加载web.xml文件，读取listener和context-param两个节点
1. 创建一个ServletContext（Servlet上下文）这个上下文供所有部分共享
1. 容器将context-param转换成键值对，交给ServletContext
1. 接着按照上述顺序继续执行

> listener是实现了javax.servlet.ServletContextListener 接口的服务器端程序, 它也是随web应用的启动
而启动，只初始化一次，随web应用的停止而销毁, 实际上就是监听 Web 应用的生命周期。主要作用是：做一些初始化的内容添加工作、设置一些基本的内容、比如一些参数或者是一些 固定的对象等等. listener被加载的顺序就是它们在web.xml中定义的顺序.

在Web容器中使用Spring MVC，就要进行四个方面的配置:

1. 编写”(servlet-name)”-servlet.xml:这里的servlet-name是标签<servlet-name>指定的值，必须是相同的，下面例子中是springmvc-servlet.xml

    ```xml
    <beans>
        <!-- 扫描包 spring可以自动去扫描base-pack下面或者子包下面的java文件，如果扫描到有@Component @Controller@Service等这些注解的类，则把这些类注册为bean-->
        <context:component-scan base-package="com.controller"/>

        <!-- don't handle the static resource -->
        <mvc:default-servlet-handler />

        <!-- 注解驱动-->
        <mvc:annotation-driven />

       <!-- 对转向页面的路径解析. prefix:前缀， suffix:后缀   如:http://127.0.0.1:8080/springmvc/jsp/****.jsp-->
        <bean class="org.springframework.web.servlet.view.InternalResourceViewResolver" 
                id="internalResourceViewResolver">
            <!-- 前缀 -->
            <property name="prefix" value="/WEB-INF/jsp/" />
            <!-- 后缀 -->
            <property name="suffix" value=".jsp" />
        </bean>
    </beans>
    ```

1. 添加servlet定义配置DispatcherServlet:前端处理器控制器，接受HTTP请求和转发请求的类，是分发Controller请求的，是Spring的核心要素

    ```xml
     <!-- 配置前端控制器DispatcherServlet -->
     <servlet>
          <servlet-name>springmvc</servlet-name>
          <servlet-class>org.springframework.web.servlet.DispatcherServlet</servlet-class>
          <!-- 初始applicationContext.xml:applicationContext.xml配置文件也可以使用<init-param>标签在servlet标签中进行配置 -->
          <init-param>
          <!-- 配置spring文件 -->
          <param-name>contextConfigLocation</param-name>
          <param-value>classpath:springmvc-servlet.xml</param-value>
          </init-param>
          <!--标记容器启动的时候就启动这个servlet-->
          <load-on-startup>1</load-on-startup>
      </servlet>

      <!-- 配置请求地址拦截url -->
      <servlet-mapping>
          <servlet-name>springmvc</servlet-name>
          <!--拦截所有-->
          <url-pattern>/</url-pattern>
      </servlet-mapping>
    ```

1. 配置contextConfigLocation初始化参数:**指定Spring IOC容器需要读取的定义了非web层的Bean（DAO/Service）的XML文件路径。可以指定多个XML文件路径，可以用逗号、冒号等来分隔**。如果没有指定”contextConfigLocation”参数，则会在 /WEB-INF/下查找 “servlet-name(就是下图中必须相同的servlet-name)-servlet.xml” 这样的文件加载，也就是springmvc-servlet.xml

    ```xml
    <!-- 如果不配置contextConfigLocation，则会默认寻找<servlet-name>标签中定义的值，也就是默认找到WEB-INF(classpath)/springmvc-servlet.xml -->
    <context-param>
        <!-- 指定spring bean的配置文件所在目录 -->
        <param-name>contextConfigLocation</param-name>
        <param-value>classpath:springmvc-servlet.xml</param-value>
    </context-param>
    ```

    > 其实`<context-param>`就是用于创建spring的 xxxApplicationContext, 比如`org.springframework.context.support.ClassPathXmlApplicationContext`

    > ClassPathXmlApplicationContext是该容器从XML 文件(spring配置)中加载已被定义的bean, 即实现了包含了BeanFactory所提供的功能

1. 配置ContextLoaderListerner:Spring MVC在Web容器中的启动类，负责Spring IOC(IOC介绍)容器在Web上下文中的初始化

    ```xml
     <listener>
        <listener-class>
            org.springframework.web.context.ContextLoaderListener
        </listener-class>
     </listener>
    ```

    ContextLoaderListener(listener-class)的作用就是启动Web容器时，自动装配ApplicationContext的配置信息. 因为它实现了ServletContextListener这个接口，在web.xml配置这个监听器，启动容器时，就会默认执行它实现的contextInitialized方法

### bean xml
Spring框架的本质其实是:通过XML配置来驱动Java代码，这样就可以把原本由java代码管理的耦合关系，提取到XML配置文件中管理. 这样就实现了系统中各组件的解耦，有利于后期的升级和维护.

beans是Spring配置文件的根元素，该元素可以指定如下属性:
- default-lazy-init:指定元素下配置的所有bean默认的延迟初始化行为
- default-merge:指定元素下配置的所有bean默认的merge行为
- default-autowire:指定元素下配置的所有bean默认的自动装配行为
- default-init-method:指定元素下配置的所有bean默认的初始化方法
- default-destroy-method:指定元素下配置的所有bean默认的回收方法
- default-autowire-candidates:指定元素下配置的所有bean默认是否作为自动装配的候选Bean

使用bean的init-method和destroy-method属性可初始化和销毁单独的bean

**bean可不实现bean xml中定义的`default-xxx`方法**.

example:
```xml
<bean id="AuthorizationManager" class="org.zstack.identity.AuthorizationManager">
        <zstack:plugin>
            <zstack:extension interface="org.zstack.header.Component"/>
            <zstack:extension interface="org.zstack.header.apimediator.GlobalApiMessageInterceptor"/>
            <zstack:extension interface="org.zstack.header.zql.ZQLQueryExtensionPoint"/>
        </zstack:plugin>
    </bean>
```
属性:
- id : bean在spring容器中的唯一id
- name : 等同id
- class : bean的实现类
- depends-on : 表示一个Bean的实例化依靠另一个Bean先实例化

### [`<bean class="org.springframework.beans.factory.config.PropertyPlaceholderConfigurer>`](https://blog.csdn.net/qyf_5445/article/details/8211136)
通过可将bean.xml的设定(bean的`<property>`)动态覆盖到`.properties`文件中(类似于同时使用多个ini), 而`.properties`文件可以作为自定义需求动态设定bean.

### beanRefContext.xml(位于Classpath的根目录下)
用来创建这个ApplicationContext实例, 配置中指明创建这个ApplicationContext所需的配置文件.

### 清理maven cache
`rm -rf ~/.m2/repository/*`

### maven添加mirror
ref:
- [aliyunmaven](https://developer.aliyun.com/mvn/guide)

```bash
mvn --version # 获取MAVEN_HOME
vim /usr/share/maven/conf/settings.xml # 找到 mirrors 标签，添加如下内容
<mirror>
  <id>aliyunmaven</id>
  <mirrorOf>*</mirrorOf>
  <name>阿里云公共仓库</name>
  <url>https://maven.aliyun.com/repository/public</url>
</mirror>
```

### spring xml中的占位符
ref:
- [Spring 4 官方文档学习（五）核心技术之SpEL](https://www.cnblogs.com/larryzeal/p/5910149.html)

```xml
    <bean id="RESTFacade" class="org.zstack.core.rest.RESTFacadeImpl">
        <property name="hostname" value="${RESTFacade.hostname:AUTO}" />
        <property name="port" value="${RESTFacade.port:8080}" />
        <property name="path" value="${RESTFacade.path:zstack}" />
    </bean>
```

RESTFacade.hostname其实是`zstack.properties`中的`RESTFacade.hostname`, 其"AUTO"表示默认值

### spring xml配置
- `import resource="applicationContext-tx.xml"/>` : include其他配置

### Reflections 的作用
Reflections通过扫描classpath，索引元数据，并且允许在运行时查询这些元数据。

获取某个类型的所有子类；比如，有一个父类是TestInterface，可以获取到TestInterface的所有子类
获取某个注解的所有类型/字段变量，支持注解参数匹配。
使用正则表达式获取所有匹配的资源文件
获取特定签名方法

```java
public class ReflectionTest {
 public static void main(String[] args) {
  // 扫包
  Reflections reflections = new Reflections(new ConfigurationBuilder()
    .forPackages("com.boothsun.reflections") // 指定路径URL
    .addScanners(new SubTypesScanner()) // 添加子类扫描工具
    .addScanners(new FieldAnnotationsScanner()) // 添加 属性注解扫描工具
    .addScanners(new MethodAnnotationsScanner() ) // 添加 方法注解扫描工具
    .addScanners(new MethodParameterScanner() ) // 添加方法参数扫描工具
    );

  // 反射出子类
  Set<Class<? extends ISayHello>> set = reflections.getSubTypesOf( ISayHello.class ) ;
  System.out.println("getSubTypesOf:" + set);

  // 反射出带有指定注解的类
  Set<Class<?>> ss = reflections.getTypesAnnotatedWith( MyAnnotation.class );
  System.out.println("getTypesAnnotatedWith:" + ss);

  // 获取带有特定注解对应的方法
  Set<Method> methods = reflections.getMethodsAnnotatedWith( MyMethodAnnotation.class ) ;
  System.out.println("getMethodsAnnotatedWith:" + methods);

  // 获取带有特定注解对应的字段
  Set<Field> fields = reflections.getFieldsAnnotatedWith( Autowired.class ) ;
  System.out.println("getFieldsAnnotatedWith:" + fields);

  // 获取特定参数对应的方法
  Set<Method> someMethods = reflections.getMethodsMatchParams(long.class, int.class);
  System.out.println("getMethodsMatchParams:" + someMethods);

  Set<Method> voidMethods = reflections.getMethodsReturn(void.class);
  System.out.println( "getMethodsReturn:" + voidMethods);

  Set<Method> pathParamMethods =reflections.getMethodsWithAnyParamAnnotated( PathParam.class);
  System.out.println("getMethodsWithAnyParamAnnotated:" + pathParamMethods);
 }
}

// method.isAnnotationPresent(GetPayContent.class) // 判断带自定义注解 GetPayContent 的 method
```

### JPA的persistence.xml
persistence.xml文件必须定义在classpath路径下的META-INF文件夹中.

JPA是Java Persistence API的简称，中文名Java持久层API，是JDK 5.0注解或XML描述对象－关系表的映射关系，并将运行期的实体对象持久化到数据库中, 通常与hibernate结合使用.

JPA包括以下三个方面的技术：
1. ORM映射元数据

    JPA支持XML和JDK注解两种元数据的形式，元数据描述对象和表之间的映射关系，框架据此景实体对象持久化到数据库表中

1. API

    用来操作实体对象，执行CRUD操作，框架在后台代替我们完成所有的事情，开发者能够从繁琐的JDBC和SQL代码中解放出来

1. 查询语言

    通过面向对象数据库的查询语言查询数据，避免程序的SQL语句的高度耦合

由org.springframework.orm.jpa.persistenceunit.PersistenceUnitReader载入.

标签:
- mapping-file : 指定映射文件的位置即声明orm对象位置的xml配置

# log4j2

## log4j2配置
ref:
- [Log4j2进阶使用(Pattern Layout详细设置)](https://www.jianshu.com/p/37ef7bc6d6eb)

`/usr/local/zstack/apache-tomcat/webapps/zstack/WEB-INF/classes/log4j2.xml`:
```xml
 <?xml version="1.0" encoding="UTF-8" standalone="no"?><Configuration monitorInterval="30" status="warn">
    <Appenders>
        <RollingFile fileName="${sys:catalina.home}/logs/management-server.log" filePattern="${sys:catalina.home}/logs/management-server-%d{yyyy-MM-dd}-%i.log.gz" name="RollingFile">
            <Masksensitivepatternlayout>
                <pattern>%d{yyyy-MM-dd HH:mm:ss,SSS} %-5p [%c{}] %X{api,task} (%t) %m%n</pattern>
                <replaces>
                   <!-- <replace regex='(password|remote_pass|dstPassword|consolePassword|remotePass)(\\":\\"|":"|": "|:)(.*?)(\\"|"| )' replacement="$1$2*****$4" />-->
                    <replace regex="(ansible_ssh_pass)(=)(.*?)( )" replacement="$1$2*****$4"/>
                </replaces>
            </Masksensitivepatternlayout>
            <Policies>
                <SizeBasedTriggeringPolicy size="450 MB"/>
            </Policies>
            <DefaultRolloverStrategy max="50">
                <Delete basePath="${sys:catalina.home}/logs/" maxDepth="1">
                    <IfFileName glob="management-server-*.log.gz">
                        <IfAny>
                        </IfAny>
                    </IfFileName>
                </Delete>
            </DefaultRolloverStrategy>
        </RollingFile>

        <RollingFile fileName="${sys:catalina.home}/logs/zstack-api.log" filePattern="${sys:catalina.home}/logs/zstack-api-%d{yyyy-MM-dd}-%i.log.gz" name="ApiRequestLogger">
            <Masksensitivepatternlayout>
                <pattern>%d{yyyy-MM-dd HH:mm:ss,SSS} %-5p [%c{}] (%t) %m%n</pattern>
                <replaces>
                    <!--<replace regex='(password|remote_pass|dstPassword|consolePassword|remotePass)(\\":\\"|":"|": "|:)(.*?)(\\"|"| )' replacement="$1$2*****$4" />-->
                    <replace regex="(ansible_ssh_pass)(=)(.*?)( )" replacement="$1$2*****$4"/>
                </replaces>
            </Masksensitivepatternlayout>
            <Policies>
                <SizeBasedTriggeringPolicy size="150 MB"/>
            </Policies>
            <DefaultRolloverStrategy max="50">
                <Delete basePath="${sys:catalina.home}/logs/" maxDepth="1">
                    <IfFileName glob="zstack-api-*.log.gz">
                        <IfAny>
                        </IfAny>
                    </IfFileName>
                </Delete>
            </DefaultRolloverStrategy>
        </RollingFile>

        <RollingFile fileName="${sys:catalina.home}/logs/zstack-disk-capacity.log" filePattern="${sys:catalina.home}/logs/zstack-disk-capacity-%d{yyyy-MM-dd}-%i.log.gz" name="DiskCapacityLogger">
            <Masksensitivepatternlayout>
                <pattern>%d{yyyy-MM-dd HH:mm:ss,SSS} %-5p [%c{}] (%t) %m%n</pattern>
                <replaces>
                   <!-- <replace regex='(password|remote_pass|dstPassword|consolePassword|remotePass)(\\":\\"|":"|": "|:)(.*?)(\\"|"| )' replacement="$1$2*****$4" />-->
                    <replace regex="(ansible_ssh_pass)(=)(.*?)( )" replacement="$1$2*****$4"/>
                </replaces>
            </Masksensitivepatternlayout>
            <Policies>
                <SizeBasedTriggeringPolicy size="150 MB"/>
            </Policies>
            <DefaultRolloverStrategy max="50"/>
        </RollingFile>

        <RollingFile fileName="${sys:catalina.home}/logs/zstack-disk-capacity-details.log" filePattern="${sys:catalina.home}/logs/zstack-disk-capacity-details-%d{yyyy-MM-dd}-%i.log.gz" name="DiskCapacityLoggerDetails">
            <Masksensitivepatternlayout>
                <pattern>%d{yyyy-MM-dd HH:mm:ss,SSS} %-5p [%c{}] (%t) %m%n</pattern>
                <replaces>
                   <!-- <replace regex='(password|remote_pass|dstPassword|consolePassword|remotePass)(\\":\\"|":"|": "| )(.*?)(\\"|"| )' replacement="$1$2*****$4" />-->
                    <replace regex="(ansible_ssh_pass)(=)(.*?)( )" replacement="$1$2*****$4"/>
                </replaces>
            </Masksensitivepatternlayout>
            <Policies>
                <SizeBasedTriggeringPolicy size="150 MB"/>
            </Policies>
            <DefaultRolloverStrategy max="50"/>
        </RollingFile>

        <Async bufferSize="512" ignoreExceptions="false" name="Async">
            <AppenderRef ref="RollingFile"/>
        </Async>
        <Rewrite name="Rewrite">
            <MaskSensitiveInfoRewritePolicy/>
            <AppenderRef ref="Async"/>
        </Rewrite>
    </Appenders>

    <Loggers>
        <Logger additivity="TRUE" level="trace" name="org.zstack.storage.primary.DiskCapacityTracer">
            <AppenderRef level="trace" ref="DiskCapacityLogger"/>
        </Logger>

        <Logger additivity="TRUE" level="trace" name="org.zstack.storage.primary.DiskCapacityTracerDetails">
            <AppenderRef level="trace" ref="DiskCapacityLoggerDetails"/>
        </Logger>

        <Logger additivity="TRUE" level="TRACE" name="api.request">
            <AppenderRef level="TRACE" ref="ApiRequestLogger"/>
        </Logger>

        <Logger level="trace" name="org.zstack"/>

        <Logger level="trace" name="org.zstack.utils"/>
        <Logger level="trace" name="org.zstack.utils.string.StringSimilarity"/>
        <Logger level="trace" name="org.zstack.utils.HTTP"/>

        <Logger level="trace" name="org.zstack.core.rest"/>

        <Logger level="trace" name="org.zstack.core.cloudbus"/>

        <Logger level="trace" name="org.springframework"/>

        <Logger level="TRACE" name="org.zstack.core.workflow"/>

        <!--
                <Logger name="org.zstack.billing" level="TRACE" />
        -->

        <Logger level="trace" name="org.zstack.zwatch.prometheus.ProgressMonitorHelper"/>

        <Logger level="trace" name="org.zstack.header.rest.TimeoutRestTemplate"/>

        <Logger level="trace" name="org.hibernate"/>

        <Logger level="TRACE" name="org.zstack.storage.primary"/>
        <!--
          Root节点用来指定项目的根日志，如果没有单独指定Logger，那么就会默认使用该Root日志输出.
          level: All < Trace < Debug < Info < Warn < Error < Fatal < OFF.
          AppenderRef：用来指定该日志输出到哪个Appender.
        -->
        <Root additivity="false" level="trace">
            <AppenderRef ref="Rewrite"/>
        </Root>
    </Loggers>
</Configuration>
```

### SingularAttribute
ref:
- [JPA criteria 查询:类型安全与面向对象](https://blog.51cto.com/u_9058648/3563753)

静态属性SingularAttribute<A, B> b，这里b是定义在类A中的类型为B的一个对象