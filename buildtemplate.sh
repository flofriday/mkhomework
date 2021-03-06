# A simple script to build the template
mkhomework \
-subject Math \
-title "Homework #1" \
-author flofriday \
-duedate 10.09.2020 \
-tasks "Problem A, Problem B, Problem C, Additional Task" \
-output examples/math-de.tex \
math-de.tmp

# A simple script to build the template
mkhomework \
-subject "Statistics and Probability" \
-title "HW #2" \
-author flofriday \
-duedate 06.10.2020 \
-tasks "1) Bonferroni's inequality, 2, 3, 4, 5, 6" \
-output examples/ws.tex \
ws.tmp

# The artificial inteligent homework
mkhomework \
-subject "Einführung in Künstliche Intelligenz" \
-title "HW #1" \
-author flofriday \
-duedate 06.10.2020 \
-tasks "Exercise 1.1, Exercise 1.2, Exercise 1.3" \
-output examples/ai.tex \
ai.tmp
 