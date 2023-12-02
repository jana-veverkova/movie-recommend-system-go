# Movie Recommendation System in Go

---
title: "Movie Recommendation System"
author: "Jana Veverková"
output: github_document
---

This recommendation system is based on edx data science course. The Movielens dataset of 10M movies and ratings was used. Based on edx course the dataset was split to an edx dataset and a test dataset. Further, only the edx dataset was used which was further split 4:1 to the train set and the test set. Parameters of the models were computed using the train set only and then, the models were evaluated using the test set. RMSE and MAE were calculated.

Data containing director and 2 main actors were web-scraped for every movie from csfd.cz.

## Modelv0 
###(0 variables)

For comparison purposes, the first model developed is the most simplest one:

```math
$$y_{ij} = a$$
```
where $a$ is the average of all ratings.

Summary statistics of this model based on test set are:

```math
$$RMSE = $$
$$MAE = $$
```

## Modelv2 
###(2 variables)

```math
$$y_{ij} = a + m_{i} + u_{j} $$
```
where $m_{ij}$ is an effect of movie $i$ calculated as an average of $y_{ij} - a$ for every $j$ and 
$m_{ij}$ is an effect of user $j$ calculated as an average of $y_{ij} - a - m_{ij}$ for every $i$

Summary statistics of this model based on test set are:

```math
$$RMSE = $$
$$MAE = $$
```

## Modelv2B
###(2 variables with lambda coeffient)

Model 
```math
$$y_{ij} = a + m_{i} + u_{j} $$
```
was adjusted so that estimates of movies and users effects were averages adjusted by lambda coefficients (sums were divided by number of observations + lambda).




