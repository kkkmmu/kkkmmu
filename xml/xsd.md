1.XML Schema 描述了 XML文档的结构。
2.<schema> 元素是每一个 XML Schema 的根元素。
	<?xml version="1.0"?>

	<xs:schema>
	...
	...
	</xs:schema>
	<schema> 元素可包含属性。一个 schema 声明往往看上去类似这样：

	<?xml version="1.0"?>

	<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"
	targetNamespace="http://www.runoob.com"
	xmlns="http://www.runoob.com"
	elementFormDefault="qualified">
	...
	...
	</xs:schema>
	
	以下代码片段:

	xmlns:xs="http://www.w3.org/2001/XMLSchema"
	显示 schema 中用到的元素和数据类型来自命名空间 "http://www.w3.org/2001/XMLSchema"。同时它还规定了来自命名空间 "http://www.w3.org/2001/XMLSchema" 的元素和数据类型应该使用前缀 xs：

	这个片断：

	targetNamespace="http://www.runoob.com"
	显示被此 schema 定义的元素 (note, to, from, heading, body) 来自命名空间： "http://www.runoob.com"。

	这个片断：

	xmlns="http://www.runoob.com"
	指出默认的命名空间是 "http://www.runoob.com"。

	这个片断：

	elementFormDefault="qualified"
	指出任何 XML 实例文档所使用的且在此 schema 中声明过的元素必须被命名空间限定。

3. 在 XML 文档中引用 Schema
	<?xml version="1.0"?>

	<note xmlns="http://www.runoob.com"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.runoob.com note.xsd">

		<to>Tove</to>
		<from>Jani</from>
		<heading>Reminder</heading>
		<body>Don''t forget me this weekend!</body>
	</note>
	下面的代码片断：

	xmlns="http://www.runoob.com"
	规定了默认命名空间的声明。此声明会告知 schema 验证器，在此 XML 文档中使用的所有元素都被声明于 "http://www.runoob.com" 这个命名空间。

	一旦您拥有了可用的 XML Schema 实例命名空间：

	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	您就可以使用 schemaLocation 属性了。此属性有两个值。第一个值是需要使用的命名空间。第二个值是供命名空间使用的 XML schema 的位置：

	xsi:schemaLocation="http://www.runoob.com note.xsd"

4. XML Schema 可定义 XML 文件的元素。
	简易元素指那些只包含文本的元素。它不会包含任何其他的元素或属性。
	不过，"仅包含文本"这个限定却很容易造成误解。文本有很多类型。它可以是 XML Schema 定义中包括的类型中的一种（布尔、字符串、数据等等），或者它也可以是您自行定义的定制类型。

	您也可向数据类型添加限定（即 facets），以此来限制它的内容，或者您可以要求数据匹配某种特定的模式。

5. 定义简易元素
	定义简易元素的语法：

	<xs:element name="xxx" type="yyy"/>
	此处 xxx 指元素的名称，yyy 指元素的数据类型。XML Schema 拥有很多内建的数据类型。

	最常用的类型是：
	xs:string
	xs:decimal
	xs:integer
	xs:boolean
	xs:date
	xs:time
	实例
	这是一些 XML 元素：

	<lastname>Refsnes</lastname>
	<age>36</age>
	<dateborn>1970-03-27</dateborn>
	这是相应的简易元素定义：

	<xs:element name="lastname" type="xs:string"/>
	<xs:element name="age" type="xs:integer"/>
	<xs:element name="dateborn" type="xs:date"/>

	简易元素的默认值和固定值
	简易元素可拥有指定的默认值或固定值。

	当没有其他的值被规定时，默认值就会自动分配给元素。

	在下面的例子中，缺省值是 "red"：

	<xs:element name="color" type="xs:string" default="red"/>
	固定值同样会自动分配给元素，并且您无法规定另外一个值。

	在下面的例子中，固定值是 "red"：

	<xs:element name="color" type="xs:string" fixed="red"/>

6.简易元素无法拥有属性。假如某个元素拥有属性，它就会被当作某种复合类型。但是属性本身总是作为简易类型被声明的。
	定义属性的语法是

	<xs:attribute name="xxx" type="yyy"/>
	在此处，xxx 指属性名称，yyy 则规定属性的数据类型。XML Schema 拥有很多内建的数据类型。

	最常用的类型是：
	xs:string
	xs:decimal
	xs:integer
	xs:boolean
	xs:date
	xs:time
	实例
	这是带有属性的 XML 元素：

	<lastname lang="EN">Smith</lastname>
	这是对应的属性定义：

	<xs:attribute name="lang" type="xs:string"/>

	属性的默认值和固定值
	属性可拥有指定的默认值或固定值。

	当没有其他的值被规定时，默认值就会自动分配给元素。

	在下面的例子中，缺省值是 "EN"：

	<xs:attribute name="lang" type="xs:string" default="EN"/>
	固定值同样会自动分配给元素，并且您无法规定另外的值。

	在下面的例子中，固定值是 "EN"：

	<xs:attribute name="lang" type="xs:string" fixed="EN"/>

	可选的和必需的属性
	在缺省的情况下，属性是可选的。如需规定属性为必选，请使用 "use" 属性：

	<xs:attribute name="lang" type="xs:string" use="required"/>

	对内容的限定
	当 XML 元素或属性拥有被定义的数据类型时，就会向元素或属性的内容添加限定。
	假如 XML 元素的类型是 "xs:date"，而其包含的内容是类似 "Hello World" 的字符串，元素将不会（通过）验证。
	通过 XML schema，您也可向您的 XML 元素及属性添加自己的限定。这些限定被称为 facet（编者注：意为(多面体的)面，可译为限定面）。

7.限定（restriction）用于为 XML 元素或者属性定义可接受的值。对 XML 元素的限定被称为 facet。
	对值的限定
	下面的例子定义了带有一个限定且名为 "age" 的元素。age 的值不能低于 0 或者高于 120：

	<xs:element name="age">
		<xs:simpleType>
			<xs:restriction base="xs:integer">
				<xs:minInclusive value="0"/>
				<xs:maxInclusive value="120"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>

	对一组值的限定
	如需把 XML 元素的内容限制为一组可接受的值，我们要使用枚举约束（enumeration constraint）。

	下面的例子定义了带有一个限定的名为 "car" 的元素。可接受的值只有：Audi, Golf, BMW：

	<xs:element name="car">
  		<xs:simpleType>
    		<xs:restriction base="xs:string">
      			<xs:enumeration value="Audi"/>
      			<xs:enumeration value="Golf"/>
      			<xs:enumeration value="BMW"/>
    		</xs:restriction>
  		</xs:simpleType>
	</xs:element>
	上面的例子也可以被写为：

	<xs:element name="car" type="carType"/>
		<xs:simpleType name="carType">
		  <xs:restriction base="xs:string">
    	  <xs:enumeration value="Audi"/>
          <xs:enumeration value="Golf"/>
          <xs:enumeration value="BMW"/>
        </xs:restriction>
	</xs:simpleType>
	注意： 在这种情况下，类型 "carType" 可被其他元素使用，因为它不是 "car" 元素的组成部分。

	对一系列值的限定
	如需把 XML 元素的内容限制定义为一系列可使用的数字或字母，我们要使用模式约束（pattern constraint）。

	下面的例子定义了带有一个限定的名为 "letter" 的元素。可接受的值只有小写字母 a - z 其中的一个：

	<xs:element name="letter">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:pattern value="[a-z]"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	下一个例子定义了带有一个限定的名为 "initials" 的元素。可接受的值是大写字母 A - Z 其中的三个：
	<xs:element name="initials">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:pattern value="[A-Z][A-Z][A-Z]"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	下一个例子也定义了带有一个限定的名为 "initials" 的元素。可接受的值是大写或小写字母 a - z 其中的三个：

	<xs:element name="initials">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:pattern value="[a-zA-Z][a-zA-Z][a-zA-Z]"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	下一个例子定义了带有一个限定的名为 choice 的元素。可接受的值是字母 x, y 或 z 中的一个：

	<xs:element name="choice">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:pattern value="[xyz]"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	下一个例子定义了带有一个限定的名为 "prodid" 的元素。可接受的值是五个阿拉伯数字的一个序列，且每个数字的范围是 0-9：

	<xs:element name="prodid">
		<xs:simpleType>
			<xs:restriction base="xs:integer">
				<xs:pattern value="[0-9][0-9][0-9][0-9][0-9]"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>

	对一系列值的其他限定
	下面的例子定义了带有一个限定的名为 "letter" 的元素。可接受的值是 a - z 中零个或多个字母：
	<xs:element name="letter">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:pattern value="([a-z])*"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	下面的例子定义了带有一个限定的名为 "letter" 的元素。可接受的值是一对或多对字母，每对字母由一个小写字母后跟一个大写字母组成。举个例子，"sToP"将会通过这种模式的验证，但是 "Stop"、"STOP" 或者 "stop" 无法通过验证：
	<xs:element name="letter">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:pattern value="([a-z][A-Z])+"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	下面的例子定义了带有一个限定的名为 "gender" 的元素。可接受的值是 male 或者 female：
	<xs:element name="gender">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:pattern value="male|female"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	下面的例子定义了带有一个限定的名为 "password" 的元素。可接受的值是由 8 个字符组成的一行字符，这些字符必须是大写或小写字母 a - z 亦或数字 0 - 9：
	<xs:element name="password">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:pattern value="[a-zA-Z0-9]{8}"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>

	对空白字符的限定
	如需规定对空白字符（whitespace characters）的处理方式，我们需要使用 whiteSpace 限定。
	下面的例子定义了带有一个限定的名为 "address" 的元素。这个 whiteSpace 限定被设置为 "preserve"，这意味着 XML 处理器不会移除任何空白字符：

	<xs:element name="address">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:whiteSpace value="preserve"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	这个例子也定义了带有一个限定的名为 "address" 的元素。这个 whiteSpace 限定被设置为 "replace"，这意味着 XML 处理器将移除所有空白字符（换行、回车、空格以及制表符）：
	<xs:element name="address">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:whiteSpace value="replace"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	这个例子也定义了带有一个限定的名为 "address" 的元素。这个 whiteSpace 限定被设置为 "collapse"，这意味着 XML 处理器将移除所有空白字符（换行、回车、空格以及制表符会被替换为空格，开头和结尾的空格会被移除，而多个连续的空格会被缩减为一个单一的空格）：
	<xs:element name="address">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:whiteSpace value="collapse"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>

	对长度的限定
	如需限制元素中值的长度，我们需要使用 length、maxLength 以及 minLength 限定。
	本例定义了带有一个限定且名为 "password" 的元素。其值必须精确到 8 个字符：
	<xs:element name="password">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:length value="8"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>
	这个例子也定义了带有一个限定的名为 "password" 的元素。其值最小为 5 个字符，最大为 8 个字符：
	<xs:element name="password">
		<xs:simpleType>
			<xs:restriction base="xs:string">
				<xs:minLength value="5"/>
				<xs:maxLength value="8"/>
			</xs:restriction>
		</xs:simpleType>
	</xs:element>

	数据类型的限定
	限定	描述
	enumeration	定义可接受值的一个列表
	fractionDigits	定义所允许的最大的小数位数。必须大于等于0。
	length	定义所允许的字符或者列表项目的精确数目。必须大于或等于0。
	maxExclusive	定义数值的上限。所允许的值必须小于此值。
	maxInclusive	定义数值的上限。所允许的值必须小于或等于此值。
	maxLength	定义所允许的字符或者列表项目的最大数目。必须大于或等于0。
	minExclusive	定义数值的下限。所允许的值必需大于此值。
	minInclusive	定义数值的下限。所允许的值必需大于或等于此值。
	minLength	定义所允许的字符或者列表项目的最小数目。必须大于或等于0。
	pattern	定义可接受的字符的精确序列。
	totalDigits	定义所允许的阿拉伯数字的精确位数。必须大于0。
	whiteSpace	定义空白字符（换行、回车、空格以及制表符）的处理方式。

8.复合元素包含了其他的元素及/或属性。
	什么是复合元素？
	复合元素指包含其他元素及/或属性的 XML 元素。
	有四种类型的复合元素：
		空元素
		包含其他元素的元素
		仅包含文本的元素
		包含元素和文本的元素
	注意： 上述元素均可包含属性！

9.复合元素的例子
	复合元素，"product"，是空的：
		<product pid="1345"/>
	复合元素，"employee"，仅包含其他元素：
	<employee>
  		<firstname>John</firstname>
  		<lastname>Smith</lastname>
	</employee>
	复合 XML 元素，"food"，仅包含文本：
	<food type="dessert">Ice cream</food>
	复合XML元素，"description"包含元素和文本：
	<description>
		It happened on <date lang="norwegian">03.03.99</date> ....
	</description>

10. 如何定义复合元素？
	请看这个复合 XML 元素，"employee"，仅包含其他元素：
	<employee>
  		<firstname>John</firstname>
  		<lastname>Smith</lastname>
	</employee>
	在 XML Schema 中，我们有两种方式来定义复合元素：
	1. 通过命名此元素，可直接对"employee"元素进行声明，就像这样：
	<xs:element name="employee">
  		<xs:complexType>
    		<xs:sequence>
      			<xs:element name="firstname" type="xs:string"/>
      			<xs:element name="lastname" type="xs:string"/>
    		</xs:sequence>
  		</xs:complexType>
	</xs:element>
	假如您使用上面所描述的方法，那么仅有 "employee" 可使用所规定的复合类型。请注意其子元素，"firstname" 以及 "lastname"，被包围在指示器 <sequence>中。这意味着子元素必须以它们被声明的次序出现。您会在 XSD 指示器 这一节学习更多有关指示器的知识。
	2. "employee" 元素可以使用 type 属性，这个属性的作用是引用要使用的复合类型的名称：
	<xs:element name="employee" type="personinfo"/>
	<xs:complexType name="personinfo">
  		<xs:sequence>
    		<xs:element name="firstname" type="xs:string"/>
    		<xs:element name="lastname" type="xs:string"/>
  		</xs:sequence>
	</xs:complexType>
	如果您使用了上面所描述的方法，那么若干元素均可以使用相同的复合类型，比如这样：
	<xs:element name="employee" type="personinfo"/>
	<xs:element name="student" type="personinfo"/>
	<xs:element name="member" type="personinfo"/>
	<xs:complexType name="personinfo">
  		<xs:sequence>
    		<xs:element name="firstname" type="xs:string"/>
    		<xs:element name="lastname" type="xs:string"/>
  		</xs:sequence>
	</xs:complexType>
	您也可以在已有的复合元素之上以某个复合元素为基础，然后添加一些元素，就像这样：
	<xs:element name="employee" type="fullpersoninfo"/>
	<xs:complexType name="personinfo">
  		<xs:sequence>
    		<xs:element name="firstname" type="xs:string"/>
    		<xs:element name="lastname" type="xs:string"/>
  		</xs:sequence>
	</xs:complexType>
	<xs:complexType name="fullpersoninfo">
  		<xs:complexContent>
    		<xs:extension base="personinfo">
      			<xs:sequence>
        			<xs:element name="address" type="xs:string"/>
        			<xs:element name="city" type="xs:string"/>
        			<xs:element name="country" type="xs:string"/>
      			</xs:sequence>
    		</xs:extension>
  		</xs:complexContent>
	</xs:complexType>

10.带有混合内容的复合类型
	XML 元素，"letter"，含有文本以及其他元素：

	<letter>
		Dear Mr.<name>John Smith</name>.
		Your order <orderid>1032</orderid>
		will be shipped on <shipdate>2001-07-13</shipdate>.
	</letter>
	下面这个 schema 声明了这个 "letter" 元素：
	<xs:element name="letter">
		<xs:complexType mixed="true">
			<xs:sequence>
				<xs:element name="name" type="xs:string"/>
				<xs:element name="orderid" type="xs:positiveInteger"/>
				<xs:element name="shipdate" type="xs:date"/>
			</xs:sequence>
		</xs:complexType>
	</xs:element>
	注意： 为了使字符数据可以出现在 "letter" 的子元素之间，mixed 属性必须被设置为 "true"。<xs:sequence> 标签 (name、orderid 以及 shipdate ) 意味着被定义的元素必须依次出现在 "letter" 元素内部。
	我们也可以为 complexType 元素起一个名字，并让 "letter" 元素的 type 属性引用 complexType 的这个名称（通过这个方法，若干元素均可引用同一个复合类型）：
	<xs:element name="letter" type="lettertype"/>
		<xs:complexType name="lettertype" mixed="true">
			<xs:sequence>
				<xs:element name="name" type="xs:string"/>
				<xs:element name="orderid" type="xs:positiveInteger"/>
				<xs:element name="shipdate" type="xs:date"/>
			</xs:sequence>
		</xs:complexType>


11. 指示器
	有七种指示器：
		Order 指示器：
			All
			Choice
			Sequence
		Occurrence 指示器：
			maxOccurs
			minOccurs
		Group 指示器：
			Group name
			attributeGroup name
	Order 指示器
		Order 指示器用于定义元素的顺序。

		All 指示器
			<all> 指示器规定子元素可以按照任意顺序出现，且每个子元素必须只出现一次：
			<xs:element name="person">
 				 <xs:complexType>
    				<xs:all>
      					<xs:element name="firstname" type="xs:string"/>
      					<xs:element name="lastname" type="xs:string"/>
    				</xs:all>
  				</xs:complexType>
			</xs:element>
			注意： 当使用 <all> 指示器时，你可以把 <minOccurs> 设置为 0 或者 1，而只能把 <maxOccurs> 指示器设置为 1（稍后将讲解 <minOccurs> 以及 <maxOccurs>）。
		Choice 指示器
			<choice> 指示器规定可出现某个子元素或者可出现另外一个子元素（非此即彼）：
			<xs:element name="person">
 				 <xs:complexType>
    				<xs:choice>
      					<xs:element name="employee" type="employee"/>
      					<xs:element name="member" type="member"/>
    				</xs:choice>
  				</xs:complexType>
			</xs:element>
		Sequence 指示器
			<sequence> 规定子元素必须按照特定的顺序出现：
				<xs:element name="person">
   					<xs:complexType>
    					<xs:sequence>
      						<xs:element name="firstname" type="xs:string"/>
      						<xs:element name="lastname" type="xs:string"/>
    					</xs:sequence>
  					</xs:complexType>
				</xs:element>
		Occurrence 指示器
			Occurrence 指示器用于定义某个元素出现的频率。
			注意： 对于所有的 "Order" 和 "Group" 指示器（any、all、choice、sequence、group name 以及 group reference），其中的 maxOccurs 以及 minOccurs 的默认值均为 1。
			maxOccurs 指示器
				<maxOccurs> 指示器可规定某个元素可出现的最大次数：
				<xs:element name="person">
  					<xs:complexType>
   						<xs:sequence>
      						<xs:element name="full_name" type="xs:string"/>
      						<xs:element name="child_name" type="xs:string" maxOccurs="10"/>
    					</xs:sequence>
  					</xs:complexType>
				</xs:element>
				上面的例子表明，子元素 "child_name" 可在 "person" 元素中最少出现一次（其中 minOccurs 的默认值是 1），最多出现 10 次。
			minOccurs 指示器
				<minOccurs> 指示器可规定某个元素能够出现的最小次数：
				<xs:element name="person">
  					<xs:complexType>
    					<xs:sequence>
      						<xs:element name="full_name" type="xs:string"/>
      						<xs:element name="child_name" type="xs:string"
      						maxOccurs="10" minOccurs="0"/>
    					</xs:sequence>
  					</xs:complexType>
				</xs:element>
			上面的例子表明，子元素 "child_name" 可在 "person" 元素中出现最少 0 次，最多出现 10 次。
			提示：如需使某个元素的出现次数不受限制，请使用 maxOccurs="unbounded" 这个声明：
			一个实际的例子：

			名为 "Myfamily.xml" 的 XML 文件：
			<?xml version="1.0" encoding="ISO-8859-1"?>
				<persons xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
					xsi:noNamespaceSchemaLocation="family.xsd">
					<person>
 					 <full_name>Hege Refsnes</full_name>
  					 <child_name>Cecilie</child_name>
					</person>
					<person>
  					 <full_name>Tove Refsnes</full_name>
  					 <child_name>Hege</child_name>
  					 <child_name>Stale</child_name>
  					 <child_name>Jim</child_name>
  					 <child_name>Borge</child_name>
					</person>
					<person>
  					 <full_name>Stale Refsnes</full_name>
					</person>
				</persons>
			上面这个 XML 文件含有一个名为 "persons" 的根元素。在这个根元素内部，我们定义了三个 "person" 元素。每个 "person" 元素必须含有一个 "full_name" 元素，同时它可以包含多至 5 个 "child_name" 元素。

			这是schema文件"family.xsd"：

			<?xml version="1.0" encoding="ISO-8859-1"?>
				<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"
					elementFormDefault="qualified">
			<xs:element name="persons">
  				<xs:complexType>
    				<xs:sequence>
      					<xs:element name="person" maxOccurs="unbounded">
        					<xs:complexType>
          						<xs:sequence>
            						<xs:element name="full_name" type="xs:string"/>
            						<xs:element name="child_name" type="xs:string"
            							minOccurs="0" maxOccurs="5"/>
          						</xs:sequence>
        					</xs:complexType>
      					</xs:element>
    				</xs:sequence>
 			   </xs:complexType>
			</xs:element>
			</xs:schema>
			Group 指示器
				Group 指示器用于定义相关的数批元素。
				元素组
				元素组通过 group 声明进行定义：

				<xs:group name="groupname">
					...
				</xs:group>
			您必须在 group 声明内部定义一个 all、choice 或者 sequence 元素。下面这个例子定义了名为 "persongroup" 的 group，它定义了必须按照精确的顺序出现的一组元素：
			<xs:group name="persongroup">
  				<xs:sequence>
    				<xs:element name="firstname" type="xs:string"/>
    				<xs:element name="lastname" type="xs:string"/>
    				<xs:element name="birthday" type="xs:date"/>
  				</xs:sequence>
			</xs:group>
			在您把 group 定义完毕以后，就可以在另一个定义中引用它了：
			<xs:group name="persongroup">
  				<xs:sequence>
    				<xs:element name="firstname" type="xs:string"/>
    				<xs:element name="lastname" type="xs:string"/>
    				<xs:element name="birthday" type="xs:date"/>
  				</xs:sequence>
			</xs:group>
			<xs:element name="person" type="personinfo"/>
			<xs:complexType name="personinfo">
  				<xs:sequence>
    				<xs:group ref="persongroup"/>
    				<xs:element name="country" type="xs:string"/>
  				</xs:sequence>
			</xs:complexType>
		属性组
			属性组通过 attributeGroup 声明来进行定义：
			<xs:attributeGroup name="groupname">
				...
			</xs:attributeGroup>
			下面这个例子定义了名为 "personattrgroup" 的一个属性组：
			<xs:attributeGroup name="personattrgroup">
  				<xs:attribute name="firstname" type="xs:string"/>
  				<xs:attribute name="lastname" type="xs:string"/>
  				<xs:attribute name="birthday" type="xs:date"/>
			</xs:attributeGroup>
			在您已定义完毕属性组之后，就可以在另一个定义中引用它了，就像这样：

			<xs:attributeGroup name="personattrgroup">
  				<xs:attribute name="firstname" type="xs:string"/>
  				<xs:attribute name="lastname" type="xs:string"/>
  				<xs:attribute name="birthday" type="xs:date"/>
			</xs:attributeGroup>
			<xs:element name="person">
			  <xs:complexType>
    			<xs:attributeGroup ref="personattrgroup"/>
  			  </xs:complexType>
			</xs:element>
