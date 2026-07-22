SELECT count(distinct(source)) as "Number of Distinct Sources",
       count(*) as "Number of Ngrams"
from ngrams;