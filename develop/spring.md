# spring

## example
```java
package com.logicbig.example;

import org.springframework.beans.factory.support.BeanDefinitionBuilder;
import org.springframework.beans.factory.support.DefaultListableBeanFactory;

public class BeanDefinitionBuilderExample {

    public static void main (String[] args) {
        DefaultListableBeanFactory beanFactory =
                  new DefaultListableBeanFactory();

        BeanDefinitionBuilder b =
                  BeanDefinitionBuilder.rootBeanDefinition(MyBean.class)
                                       .addPropertyValue("str", "myStringValue");

        beanFactory.registerBeanDefinition("myBean", b.getBeanDefinition()); // 注册了一个新 bean


        MyBean bean = beanFactory.getBean(MyBean.class);
        bean.doSomething();
    }

    private static class MyBean {
        private String str;

        public void setStr (String str) {
            this.str = str;
        }

        public void doSomething () {
            System.out.println("from MyBean " + str);
        }
    }
}
```

`BeanDefinitionBuilder.rootBeanDefinition(MyBean.class).addPropertyValue("str", "myStringValue");`效果是定义了一个class:
```java
private static class MyBean {
    private String str = "myStringValue";

    public void doSomething () {
        System.out.println("from MyBean " + str);
    }
}
```

## bean
### 属性
- scope

    - Singleton

        这也是Spring默认的scope，表示Spring容器只创建一个bean的实例，Spring在创建第一次后会缓存起来，之后不再创建，就是设计模式中的单例模式。

    - Prototype

        代表线程每次调用这个bean都新创建一个实例。

    - Request

        表示每个request作用域内的请求只创建一个实例。

    - Session

        表示每个session作用域内的请求只创建一个实例。

    - GlobalSession

        这个只在porlet的web应用程序中才有意义，它映射到porlet的global范围的session，如果普通的web应用使用了这个scope，容器会把它作为普通的session作用域的scope创建
- name

    默认使用bean id作为name, 没有时再使用bean name