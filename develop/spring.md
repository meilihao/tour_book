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