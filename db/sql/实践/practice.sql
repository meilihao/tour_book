-- 基于data.sql的实践

-- 查询“001”课程比“002”课程成绩高的所有学生的学号 : 自连接
select
  s1.student_no
from 
  score s1 inner join score s2 on s1.student_no = s2.student_no
where 
  s1.course_no = 001 
  and 
  s2.course_no = 002 
  and
  s1.score > s2.score;

-- 查询平均成绩大于60分的同学的学号和平均成绩 : avg = sum()/行数
-- 每个学生的总分/课程数, 课程数=行数
select
  s1.student_no,
  AVG(s1.score) a
from
  score s1
  group by s1.student_no
  having a > 60 order by a;

-- 查询所有同学的学号、姓名、选课数、总成绩
select stu1.student_no, stu1.name, s.c, s.s
from
    student stu1 inner join (
        select
        s1.student_no,
        count(s1.course_no) c,
        sum(s1.score) s
        from
        score s1
        group by s1.student_no
    ) s on stu1.student_no = s.student_no;

-- 查询姓“李”的老师的个数 : count + like
select
  count(*)
from
  teacher t1
where
  t1.name like '李%'

-- 查询没学过“叶平”老师课的同学的学号、姓名
select
  stu1.student_no,
  stu1.name
from
  student stu1
where
  stu1.student_no not in
  (
    -- 学过叶平老师课的学生学号
    select distinct
        s1.student_no
    from
        score s1,
        course c1,
        teacher t1
    where
        s1.course_no = c1.course_no
        and
        c1.teacher_no = t1.teacher_no
        and
        t1.name = '叶平'
  );

select
  stu1.student_no,
  stu1.name
from
  student stu1
where
  stu1.student_no not in
  (
    -- 学过叶平老师课的学生学号
    select distinct
        s1.student_no
    from
        course c1
        inner join score s1  on s1.course_no = c1.course_no
    where
    c1.teacher_no in (select teacher_no from teacher WHERE name = '叶平')
  );

-- 查询学过“001”并且也学过编号“002”课程的同学的学号、姓名
select
  stu1.student_no,
  stu1.name
from
  score s1,
  student stu1
where
  s1.student_no = stu1.student_no
  and
  s1.course_no = 001
  and
  s1.student_no in
  (
  select distinct
    s2.student_no
  from
    score s2
  where
    s2.course_no = 002
  )

---- best
select
  stu1.student_no,
  stu1.name
from
  student stu1
where
  stu1.student_no in (
      select student_no from score where course_no in (001,002) group by student_no having count(*) = 2
  )

-- 查询学过“叶平”老师所教的所有课的同学的学号、姓名
select
  s1.student_no
from
  score s1,
  student stu1,
  course c1,
  teacher t1
where
  s1.student_no = stu1.student_no
  and
  s1.course_no = c1.course_no
  and
  c1.teacher_no = t1.teacher_no
  and
  t1.name = '叶平'
  group by s1.student_no
  having count(c1.course_no)=
  (
    select
      count(c2.course_no)
    from
      course c2,
      teacher t2
    where
      c2.teacher_no = t2.teacher_no
      and
      t2.name = '叶平'
   )

-- +------------+
-- | student_no |
-- +------------+
-- |          6 |
-- |          1 |
-- |          5 |
-- |          2 |
-- |          3 |
-- +------------+

-- 用程序拆分实现简单
-- select distinct course_no from course where teacher_no in (select teacher_no from teacher WHERE name = '叶平'); => 1,4
-- select student_no from score where course_no in (1,4) group by student_no having count(*) =2;

-- 查询有课程成绩小于60分的同学的学号、姓名
select distinct
  s1.student_no,
  stu1.name
from
  score s1,
  student stu1
where
  s1.student_no = stu1.student_no
  and
  s1.score < 60

-- 查询没有学全所有课的同学的学号、姓名
select
  s1.student_no
from
  score s1,
  student stu1
where
  s1.student_no = stu1.student_no
  group by s1.student_no
  having count(s1.course_no) <
  (
  select
    count(*)
  from
    course c1
  );

-- 查询至少有一门课与学号为“001”的同学所学相同的同学的学号和姓名
select distinct
  s1.student_no,
  stu1.name
from
  score s1,
  student stu1
where
  s1.student_no = stu1.student_no
  and
  s1.course_no in
  (
  select
    s2.course_no
  from
    score s2
  where
    s2.student_no = 001
  );